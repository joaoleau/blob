package usecases

import (
	"context"
	"github.com/joaoleau/blob/models"
	"github.com/joaoleau/blob/repository"
	"github.com/joaoleau/blob/utils"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
)

type BlobUseCase struct {
	repository repository.BlobRepo
}

func NewBlobUseCase(repo repository.BlobRepo) BlobUseCase {
	return BlobUseCase{
		repository: repo,
	}
}

// RegisterBlob registers a new blob
func (u *BlobUseCase) RegisterBlob(ctx context.Context, blob *models.Blob) (*models.Blob, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BlobUseCase.RegisterBlob")
	defer span.Finish()

	createdBlob, err := u.repository.Create(ctx, blob)
	if err != nil {
		return nil, err
	}

	return createdBlob, nil
}


func (u *BlobUseCase) UpdateBlob(ctx context.Context, blob *models.Blob) (*models.Blob, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BlobUseCase.UpdateBlob")
	defer span.Finish()

	updatedBlob, err := u.repository.Update(ctx, blob)
	if err != nil {
		return nil, err
	}

	return updatedBlob, nil
}


func (u *BlobUseCase) DeleteBlob(ctx context.Context, blobID uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BlobUseCase.DeleteBlob")
	defer span.Finish()

	if err := u.repository.Delete(ctx, blobID); err != nil {
		return err
	}

	return nil
}


func (u *BlobUseCase) GetBlobByID(ctx context.Context, blobID uuid.UUID) (*models.Blob, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BlobUseCase.GetBlobByID")
	defer span.Finish()

	blob, err := u.repository.GetByID(ctx, blobID)
	if err != nil {
		return nil, err
	}

	return blob, nil
}


func (u *BlobUseCase) ListBlobs(ctx context.Context, pq *utils.PaginationQuery) (*models.BlobList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BlobUseCase.ListBlobs")
	defer span.Finish()

	blobList, err := u.repository.ListBlobs(ctx, pq)
	if err != nil {
		return nil, err
	}

	return blobList, nil
}
