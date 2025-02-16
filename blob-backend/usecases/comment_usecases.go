package usecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/joaoleau/blob/models"
	"github.com/joaoleau/blob/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type CommentUseCase struct {
	commentRepo   repository.CommentRepo
	BlobUseCase   *BlobUseCase
}

func NewCommentUseCase(repo repository.CommentRepo, blobUseCase *BlobUseCase) CommentUseCase {
	return CommentUseCase{
		commentRepo:  repo,
		BlobUseCase: blobUseCase,
	}
}

func (c *CommentUseCase) AddComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommentUseCase.AddComment")
	defer span.Finish()

	_, err := c.BlobUseCase.GetBlobByID(ctx, comment.BlobID)
	if err != nil {
		return nil, errors.Wrap(err, "LikeUseCase.AddLike.GetByID")
	}

	email, ok := ctx.Value("email").(string)
	if !ok || email == "" {
		return nil, errors.New("user email not found in context")
	}

	user, err := c.BlobUseCase.UserUseCase.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch user by email")
	}
	if user == nil {
		return nil, errors.New("authenticated user not found")
	}

	comment.ID = uuid.New()
	comment.UserID = user.ID 
	
	newComment, err := c.commentRepo.AddComment(ctx, comment)
	if err != nil {
		return nil, errors.Wrap(err, "CommentUseCase.AddComment.AddCommentRepo")
	}

	return newComment, nil
}


func (c *CommentUseCase) RemoveComment(ctx context.Context, commentID uuid.UUID) error {

	email, ok := ctx.Value("email").(string)
	if !ok || email == "" {
		return errors.New("user email not found in context")
	}

	user, err := c.BlobUseCase.UserUseCase.GetUserByEmail(ctx, email)
	if err != nil {
		return errors.Wrap(err, "failed to fetch user by email")
	}
	if user == nil {
		return errors.New("authenticated user not found")
	}

	err = c.commentRepo.RemoveComment(ctx, user.ID, commentID)
	if err != nil {
		return errors.Wrap(err, "failed to remove comment from repository")
	}

	return nil
}

func (c *CommentUseCase) ListCommentsByBlobID(ctx context.Context, blobID uuid.UUID) ([]models.CommentWithUser, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommentUseCase.ListCommentsByBlobID")
	defer span.Finish()

	comments, err := c.commentRepo.ListCommentsByBlobID(ctx, blobID)
	if err != nil {
		return nil, errors.Wrap(err, "CommentUseCase.ListCommentsByBlobID.ListCommentsByBlobID")
	}

	return comments, nil
}