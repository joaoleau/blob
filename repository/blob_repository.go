package repository

import (
	"context"
	"database/sql"
	"github.com/joaoleau/blob/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/joaoleau/blob/models"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type BlobRepo struct {
	db *sqlx.DB
}

func NewBlobRepository(db *sqlx.DB) BlobRepo {
	return BlobRepo {db: db}
}

func (r *BlobRepo) Create(ctx context.Context, blob *models.Blob) (*models.Blob, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BlobRepo.Create")
	defer span.Finish()

	newBlob := &models.Blob{}
	if err := r.db.QueryRowxContext(ctx, createBlobQuery,
		blob.BlobID, blob.UserID, blob.Description, blob.Title,
	).StructScan(newBlob); err != nil {
		return nil, errors.Wrap(err, "BlobRepo.Create.StructScan")
	}

	return newBlob, nil
}

func (r *BlobRepo) Update(ctx context.Context, blob *models.Blob) (*models.Blob, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BlobRepo.Update")
	defer span.Finish()

	updatedBlob := &models.Blob{}
	if err := r.db.GetContext(ctx, updatedBlob, updateBlobQuery,
		blob.Description, blob.Title, blob.BlobID,
	); err != nil {
		return nil, errors.Wrap(err, "BlobRepo.Update.GetContext")
	}

	return updatedBlob, nil
}

func (r *BlobRepo) GetByID(ctx context.Context, blobID uuid.UUID) (*models.Blob, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BlobRepo.GetByID")
	defer span.Finish()

	blob := &models.Blob{}
	if err := r.db.GetContext(ctx, blob, getBlobByIDQuery, blobID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "BlobRepo.GetByID.GetContext")
	}

	return blob, nil
}

func (r *BlobRepo) Delete(ctx context.Context, blobID uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BlobRepo.Delete")
	defer span.Finish()

	_, err := r.db.ExecContext(ctx, deleteBlobQuery, blobID)
	if err != nil {
		return errors.Wrap(err, "BlobRepo.Delete.ExecContext")
	}

	return nil
}

func (r *BlobRepo) ListBlobs(ctx context.Context, pq *utils.PaginationQuery) (*models.BlobList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BlobRepo.ListBlobs")
	defer span.Finish()

	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getTotalBlob); err != nil {
		return nil, errors.Wrap(err, "BlobRepo.ListBlobs.GetContext.totalCount")
	}

	if totalCount == 0 {
		return &models.BlobList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
			Page:       pq.GetPage(),
			Size:       pq.GetSize(),
			HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
			Users:      make([]*models.Blob, 0),
		}, nil
	}

	var blobs = make([]*models.Blob, 0, pq.GetSize())
	if err := r.db.SelectContext(
		ctx,
		&blobs,
		listBlobsQuery,
		pq.GetOrderBy(),
		pq.GetOffset(),
		pq.GetLimit(),
	); err != nil {
		return nil, errors.Wrap(err, "BlobRepo.ListBlobs.SelectContext")
	}

	return &models.BlobList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
		Users:      blobs,
	}, nil
}
