package repository

import (
	"context"
	"database/sql"
	"time"
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

func (r *BlobRepo) Create(ctx context.Context, blob *models.BlobWithInterests) (*models.BlobWithInterests, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BlobRepo.Create")
	defer span.Finish()

	newBlob := &models.BlobWithInterests{}
	if err := r.db.QueryRowxContext(ctx, createBlobQuery,
		blob.ID, blob.UserID, blob.Content,
	).StructScan(newBlob); err != nil {
		return nil, errors.Wrap(err, "BlobRepo.Create.StructScan")
	}

	for _, interest := range blob.Interests {
		if err := r.createBlobInterest(ctx, newBlob.ID, interest); err != nil {
			return nil, errors.Wrap(err, "BlobRepo.Create.createBlobInterest")
		}
	}

	return newBlob, nil
}

func (r *BlobRepo) createBlobInterest(ctx context.Context, blobID uuid.UUID, interestID string) error {
	_, err := r.db.ExecContext(ctx, insertBlobInterest, blobID, interestID)
	return err
}

func (r *BlobRepo) ListAllInterests(ctx context.Context) ([]*models.Interest, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BlobRepo.ListAllInterests")
	defer span.Finish()

	var interests []*models.Interest
	query := `SELECT * FROM "Interest"`
	err := r.db.SelectContext(ctx, &interests, query)
	if err != nil {
		return nil, errors.Wrap(err, "BlobRepo.ListAllInterests.SelectContext")
	}

	return interests, nil
}


func (r *BlobRepo) GetByID(ctx context.Context, blobID uuid.UUID) (*models.BlobWithDetails, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BlobRepo.GetByID")
	defer span.Finish()

	type Row struct { 
		BlobID            string     `db:"blob_id"`
		UserID            string     `db:"blob_user_id"`
		Content           string     `db:"blob_content"`
		CreatedAt         time.Time  `db:"blob_created_at"`
		UpdatedAt         time.Time  `db:"blob_updated_at"`
		Username          string     `db:"user_username"`
		AvatarIcon        string     `db:"user_avatar_icon"`
		UserCreatedAt     time.Time  `db:"user_created_at"`
		CommentID         *string    `db:"comment_id"`
		CommentUserID     *string    `db:"comment_user_id"`
		CommentContent    *string    `db:"comment_content"`
		CommentCreatedAt  *time.Time `db:"comment_created_at"`
		CommentUpdatedAt  *time.Time `db:"comment_updated_at"`
		LikeID            *string    `db:"like_id"`
		LikeUserID        *string    `db:"like_user_id"`
		LikeCreatedAt     *time.Time `db:"like_created_at"`
		InterestID        *string    `db:"interest_id"`
		InterestName      *string    `db:"interest_name"`
		InterestDescription *string  `db:"interest_description"`
		InterestCreatedAt *time.Time `db:"interest_created_at"`
		InterestUpdatedAt *time.Time `db:"interest_updated_at"`
	}

	var rows []Row
	if err := r.db.SelectContext(ctx, &rows, getBlobByIDQuery, blobID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "BlobRepo.GetByID.SelectContext")
	}

	if len(rows) == 0 {
		return nil, nil
	}

	blob := &models.BlobWithDetails{
		ID:           rows[0].BlobID,
		UserID:       rows[0].UserID,
		Content:      rows[0].Content,
		CreatedAt:    rows[0].CreatedAt,
		UpdatedAt:    rows[0].UpdatedAt,
		Username:     rows[0].Username,
		AvatarIcon:   rows[0].AvatarIcon,
		UserCreatedAt: rows[0].UserCreatedAt,
		Comments:     []models.Comment{},
		Likes:        []models.Like{},
		Interests:    []models.Interest{},
	}

	commentIDs := make(map[string]bool)
	likeIDs := make(map[string]bool)
	interestIDs := make(map[string]bool)
	
	for _, row := range rows {
		if row.CommentID != nil {
			commentID := *row.CommentID
	
			if !commentIDs[commentID] {
				blob.Comments = append(blob.Comments, models.Comment{
					ID:        uuid.MustParse(commentID),
					Content:   *row.CommentContent,
					CreatedAt: *row.CommentCreatedAt,
					UpdatedAt: *row.CommentUpdatedAt,
					UserID:    *row.CommentUserID,
					BlobID:    uuid.MustParse(row.BlobID),
				})
				commentIDs[commentID] = true
			}
		}
	
		if row.LikeID != nil {
			likeID := *row.LikeID
			if !likeIDs[likeID] {
				blob.Likes = append(blob.Likes, models.Like{
					ID:        uuid.MustParse(likeID),
					UserID:    *row.LikeUserID,
					CreatedAt: *row.LikeCreatedAt,
					BlobID:    uuid.MustParse(row.BlobID),
				})
				likeIDs[likeID] = true
			}
		}
	
		if row.InterestID != nil {
			interestID := *row.InterestID
	
			if !interestIDs[interestID] {
				blob.Interests = append(blob.Interests, models.Interest{
					ID:          uuid.MustParse(interestID),
					Name:        *row.InterestName,
					Description: *row.InterestDescription,
					CreatedAt:   *row.InterestCreatedAt,
					UpdatedAt:   *row.InterestUpdatedAt,
				})
				interestIDs[interestID] = true
			}
		}
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

func (r *BlobRepo) ListBlobs(ctx context.Context) ([]models.BlobListWithDetails, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BlobRepo.ListBlobs")
	defer span.Finish()

	type Row struct {
		ID           string	   `db:"id"`
		UserID       string    `db:"user_id"`
		Content      string    `db:"content"`
		CreatedAt    time.Time `db:"created_at"`
		UpdatedAt    time.Time `db:"updated_at"`
		Username     string    `db:"username"`
		AvatarIcon   string    `db:"avatar_icon"`
		UserCreatedAt time.Time `db:"user_created_at"`
		InterestName *string   `db:"interest_name"`
		LikesCount   int       `db:"likes_count"`
		CommentsCount int      `db:"comments_count"`
	}

	var rows []Row
	if err := r.db.SelectContext(ctx, &rows, listBlobsQuery); err != nil {
		return nil, errors.Wrap(err, "BlobRepo.ListBlobs.SelectContext")
	}

	blobMap := make(map[string]*models.BlobListWithDetails)

	for _, row := range rows {
		if _, exists := blobMap[row.ID]; !exists {
			blobMap[row.ID] = &models.BlobListWithDetails{
				ID:          row.ID,
				UserID:      row.UserID,
				Content:     row.Content,
				CreatedAt:   row.CreatedAt,
				UpdatedAt:   row.UpdatedAt,
				Username:    row.Username,
				AvatarIcon:  row.AvatarIcon,
				UserCreatedAt: row.UserCreatedAt,
				LikesCount:  row.LikesCount,
				CommentsCount: row.CommentsCount,
				Interests:   []string{},
			}
		}

		if row.InterestName != nil {
			blobMap[row.ID].Interests = append(blobMap[row.ID].Interests, *row.InterestName)
		}
	}

	blobs := make([]models.BlobListWithDetails, 0, len(blobMap))
	for _, blob := range blobMap {
		blobs = append(blobs, *blob)
	}

	return blobs, nil
}
