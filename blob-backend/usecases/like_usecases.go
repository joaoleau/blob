package usecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/joaoleau/blob/models"
	"github.com/joaoleau/blob/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type LikeUseCase struct {
	likeRepo   repository.LikeRepo
	BlobUseCase   *BlobUseCase
}

func NewLikeUseCase(repo repository.LikeRepo, blobUseCase *BlobUseCase) LikeUseCase {
	return LikeUseCase{
		likeRepo:  repo,
		BlobUseCase: blobUseCase,
	}
}

func (l *LikeUseCase) AddLike(ctx context.Context, blobID uuid.UUID) (*models.Like, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LikeUseCase.AddLike")
	defer span.Finish()

	_, err := l.BlobUseCase.GetBlobByID(ctx, blobID)
	if err != nil {
		return nil, errors.Wrap(err, "LikeUseCase.AddLike.GetByID")
	}

	email, ok := ctx.Value("email").(string)
	if !ok || email == "" {
		return nil, errors.New("user email not found in context")
	}

	user, err := l.BlobUseCase.UserUseCase.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch user by email")
	}
	if user == nil {
		return nil, errors.New("authenticated user not found")
	}

	newLike, err := l.likeRepo.AddLike(ctx, uuid.New(), user.ID, blobID)
	if err != nil {
		return nil, errors.Wrap(err, "LikeUseCase.AddLike.AddLikeRepo")
	}

	return newLike, nil
}

func (l *LikeUseCase) RemoveLike(ctx context.Context, blobID uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LikeUseCase.RemoveLike")
	defer span.Finish()

	email, ok := ctx.Value("email").(string)
	if !ok || email == "" {
		return errors.New("user email not found in context")
	}

	user, err := l.BlobUseCase.UserUseCase.GetUserByEmail(ctx, email)
	if err != nil {
		return errors.Wrap(err, "failed to fetch user by email")
	}
	if user == nil {
		return errors.New("authenticated user not found")
	}

	likeID, err := l.likeRepo.FindLikeID(ctx, user.ID, blobID)
	if err != nil {
		return errors.Wrap(err, "Error fetching like ID")
	}
	if likeID == uuid.Nil {
		return errors.Wrap(err, "No like found for the given user and blob")
	}

	if err := l.likeRepo.RemoveLike(ctx, likeID, user.ID, blobID); err != nil {
		return errors.Wrap(err, "failed to remove like from repository")
	}

	return nil
}

func (l *LikeUseCase) ListLikesByBlobID(ctx context.Context, blobID uuid.UUID) ([]models.LikeWithUser, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LikeUseCase.ListLikesByBlobID")
	defer span.Finish()

	likes, err := l.likeRepo.ListLikesByBlobID(ctx, blobID)
	if err != nil {
		return nil, errors.Wrap(err, "LikeUseCase.ListLikesByBlobID.ListLikesByBlobIDRepo")
	}

	return likes, nil
}