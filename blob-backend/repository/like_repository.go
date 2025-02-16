package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/joaoleau/blob/models"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type LikeRepo struct {
	db *sqlx.DB
	blobRepo *BlobRepo
}

func NewLikeRepository(db *sqlx.DB, blobRepo *BlobRepo) LikeRepo {
	return LikeRepo{
		db: db, 
		blobRepo: blobRepo,
	}
}

func (r *LikeRepo) AddLike(ctx context.Context, likeID uuid.UUID, userID string, blobID uuid.UUID) (*models.Like, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LikeRepo.AddLike")
	defer span.Finish()


	if _, err := r.blobRepo.GetByID(ctx, blobID); err != nil {
		return nil, errors.Wrap(err, "LikeRepo.AddLike.GetByID")
	}

	newLike := &models.Like{}
	if err := r.db.QueryRowxContext(ctx, insertLikeQuery,
		likeID, userID, blobID,
	).StructScan(newLike); err != nil {
		return nil, errors.Wrap(err, "LikeRepo.AddLike.StructScan")
	}

	return newLike, nil
}

func (r *LikeRepo) RemoveLike(ctx context.Context, likeID uuid.UUID, userID string, blobID uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LikeRepo.RemoveLike")
	defer span.Finish()

	if _, err := r.db.ExecContext(ctx, deleteLikeQuery, likeID, userID, blobID); err != nil {
		return errors.Wrap(err, "LikeRepo.RemoveLike.ExecContext")
	}

	return nil
}


func (r *LikeRepo) FindLikeID(ctx context.Context, userID string, blobID uuid.UUID) (uuid.UUID, error) {
	var likeID uuid.UUID
	err := r.db.QueryRowContext(ctx, searchLikeQuery, userID, blobID).Scan(&likeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return uuid.Nil, nil
		}
		return uuid.Nil, errors.Wrap(err, "failed to query like ID")
	}

	return likeID, nil
}

func (r *LikeRepo) ListLikesByBlobID(ctx context.Context, blobID uuid.UUID) ([]models.LikeWithUser, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LikeRepo.ListLikesByBlobID")
	defer span.Finish()

	var likes []models.LikeWithUser
	if err := r.db.SelectContext(ctx, &likes, searchLikebyBlobIDQuery, blobID); err != nil {
		return nil, errors.Wrap(err, "LikeRepo.ListLikesByBlobID.SelectContext")
	}

	return likes, nil
}