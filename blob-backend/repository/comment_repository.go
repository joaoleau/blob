package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/joaoleau/blob/models"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/google/uuid"
)

type CommentRepo struct {
	db *sqlx.DB
	blobRepo *BlobRepo
}

func NewCommentRepository(db *sqlx.DB, blobRepo *BlobRepo) CommentRepo {
	return CommentRepo{
		db: db, 
		blobRepo: blobRepo,
	}
}

func (r *CommentRepo) AddComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommentRepo.AddComment")
	defer span.Finish()

	if _, err := r.blobRepo.GetByID(ctx, comment.BlobID); err != nil {
		return nil, errors.Wrap(err, "Comment.AddLike.GetByID")
	}

	newComment := &models.Comment{}
	if err := r.db.QueryRowxContext(ctx, insertCommentQuery,
		comment.ID, comment.Content, comment.UserID, comment.BlobID,
	).StructScan(newComment); err != nil {
		return nil, errors.Wrap(err, "CommentRepo.AddComment.StructScan")
	}

	return newComment, nil
}


func (r *CommentRepo) RemoveComment(ctx context.Context, userID string, commentID uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommentRepo.RemoveComment")
	defer span.Finish()

	if _, err := r.db.ExecContext(ctx, deleteCommentQuery, commentID, userID); err != nil {
		return errors.Wrap(err, "CommentRepo.RemoveComment.ExecContext")
	}

	return nil
}


func (r *CommentRepo) ListCommentsByBlobID(ctx context.Context, blobID uuid.UUID) ([]models.CommentWithUser, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommentRepo.ListCommentsByBlobID")
	defer span.Finish()

	var comments []models.CommentWithUser
	if err := r.db.SelectContext(ctx, &comments, searchCommentsbyBlobIDQuery, blobID); err != nil {
		return nil, errors.Wrap(err, "CommentRepo.ListCommentsByBlobID.SelectContext")
	}
	return comments, nil
}
