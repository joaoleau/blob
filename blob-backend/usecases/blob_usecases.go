package usecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/joaoleau/blob/models"
	"github.com/joaoleau/blob/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type BlobUseCase struct {
	repository  repository.BlobRepo
	UserUseCase *UserUseCase
}

func NewBlobUseCase(repo repository.BlobRepo, userUseCase *UserUseCase) BlobUseCase {
	return BlobUseCase{
		repository:  repo,
		UserUseCase: userUseCase,
	}
}

func (u *BlobUseCase) RegisterBlob(ctx context.Context, blob *models.BlobWithInterests) (*models.BlobWithInterests, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BlobUseCase.RegisterBlob")
	defer span.Finish()

	email, ok := ctx.Value("email").(string)
	if !ok || email == "" {
		return nil, errors.New("user email not found in context")
	}

	user, err := u.UserUseCase.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch user by email")
	}
	if user == nil {
		return nil, errors.New("authenticated user not found")
	}
	
	blob.ID = uuid.New()
	blob.UserID = user.ID

	createdBlob, err := u.repository.Create(ctx, blob)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create blob")
	}

	return createdBlob, nil
}

func (u *BlobUseCase) ListInterests(ctx context.Context) ([]*models.Interest, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BlobUseCase.ListInterests")
	defer span.Finish()

	interests, err := u.repository.ListAllInterests(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list interests")
	}

	return interests, nil
}


func (u *BlobUseCase) DeleteBlob(ctx context.Context, blobID uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BlobUseCase.DeleteBlob")
	defer span.Finish()

	if err := u.repository.Delete(ctx, blobID); err != nil {
		return err
	}

	return nil
}


func (u *BlobUseCase) GetBlobByID(ctx context.Context, blobID uuid.UUID) (*models.BlobWithDetails, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BlobUseCase.GetBlobByID")
	defer span.Finish()

	blobs, err := u.repository.GetByID(ctx, blobID)
	if err != nil {
		return nil, err
	}

	return blobs, nil
}


func (u *BlobUseCase) ListBlobs(ctx context.Context) ([]*models.BlobListWithDetails, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BlobUseCase.ListBlobs")
	defer span.Finish()

	blobs, err := u.repository.ListBlobs(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "BlobUseCase.ListBlobs.repository.ListBlobs")
	}

	var blobPointers []*models.BlobListWithDetails
	for _, blob := range blobs {
		blobCopy := blob
		blobPointers = append(blobPointers, &blobCopy)
	}

	if len(blobPointers) == 0 {
		return nil, nil
	}

	return blobPointers, nil
}
