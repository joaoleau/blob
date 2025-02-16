package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/joaoleau/blob/models"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*models.UserWithBlobs, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepo.GetByUsername")
	defer span.Finish()

	type Row struct {
		ID             string    `db:"id"`
		Name           string    `db:"name"`
		Email          string    `db:"email"`
		EmailVerified  time.Time `db:"email_verified"`
		Image          string    `db:"image"`
		Username       string    `db:"username"`
		Bio            string    `db:"bio"`
		AvatarIcon     string    `db:"avatar_icon"`
		AvatarColor    string    `db:"avatar_color"`
		CreatedAt      time.Time `db:"created_at"`
		UpdatedAt      time.Time `db:"updated_at"`
		BlobID         *uuid.UUID   `db:"blob_id"`
		BlobContent    *string   `db:"blob_content"`
		BlobCreatedAt  *time.Time `db:"blob_created_at"`
		BlobUpdatedAt  *time.Time `db:"blob_updated_at"`
	}

	var rows []Row

	if err := r.db.SelectContext(ctx, &rows, listUserByUsernameWithBlobsQuery, username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "UserRepo.GetByUsername.SelectContext")
	}

	if len(rows) == 0 {
		return nil, nil
	}

	userWithBlobs := &models.UserWithBlobs{
		ID:            rows[0].ID,
		Name:          rows[0].Name,
		Email:         rows[0].Email,
		EmailVerified: rows[0].EmailVerified,
		Image:         rows[0].Image,
		Username:      rows[0].Username,
		Bio:           rows[0].Bio,
		AvatarIcon:    rows[0].AvatarIcon,
		AvatarColor:   rows[0].AvatarColor,
		CreatedAt:     rows[0].CreatedAt,
		UpdatedAt:     rows[0].UpdatedAt,
		Blobs:         []models.Blob{},
	}

	for _, row := range rows {
		if row.BlobID != nil {
			blob := models.Blob{
				ID:        *row.BlobID,
				Content:   *row.BlobContent,
				CreatedAt: *row.BlobCreatedAt,
				UpdatedAt: *row.BlobUpdatedAt,
			}
			userWithBlobs.Blobs = append(userWithBlobs.Blobs, blob)
		}
	}

	return userWithBlobs, nil
}

func (r *UserRepo) GetUserById(ctx context.Context, userID string) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepo.GetUserById")
	defer span.Finish()

	user := &models.User{}
	
	if err := r.db.GetContext(ctx, user, getUserByID, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "UserRepo.GetUserById.GetContext")
	}
	return user, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepo.GetByEmail")
	defer span.Finish()

	user := &models.User{}
	
	if err := r.db.GetContext(ctx, user, getUserByEmail, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "UserRepo.GetByEmail.GetContext")
	}
	return user, nil
}

func (r *UserRepo) UpdateUser(ctx context.Context, userID string, updatedData models.User) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepo.UpdateUser")
	defer span.Finish()

	query := `UPDATE "User" SET`
	var args []interface{}
	argIndex := 1
	var oldEmail string
	
	if updatedData.Name != "" {
		query += ` name = $` + fmt.Sprintf("%d", argIndex) + `,`
		args = append(args, updatedData.Name)
		argIndex++
	}
	if updatedData.Email != "" {
		oldEmail = ctx.Value("email").(string)
		query += ` email = $` + fmt.Sprintf("%d", argIndex) + `,`
		args = append(args, updatedData.Email)
		argIndex++
	}
	if updatedData.Bio != "" {
		query += ` bio = $` + fmt.Sprintf("%d", argIndex) + `,`
		args = append(args, updatedData.Bio)
		argIndex++
	}
	if updatedData.Image != "" {
		query += ` image = $` + fmt.Sprintf("%d", argIndex) + `,`
		args = append(args, updatedData.Image)
		argIndex++
	}
	if updatedData.AvatarIcon != "" {
		query += ` avatar_icon = $` + fmt.Sprintf("%d", argIndex) + `,`
		args = append(args, updatedData.AvatarIcon)
		argIndex++
	}
	if updatedData.AvatarColor != "" {
		query += ` avatar_color = $` + fmt.Sprintf("%d", argIndex) + `,`
		args = append(args, updatedData.AvatarColor)
		argIndex++
	}

	query = query[:len(query)-1]
	query += ` WHERE id = $` + fmt.Sprintf("%d", argIndex)
	args = append(args, userID)

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "UserRepo.UpdateUser.ExecContext")
	}

	if oldEmail != updatedData.Email {
		deleteQuery := `DELETE FROM "VerificationToken" WHERE email = $1`
		_, err := r.db.ExecContext(ctx, deleteQuery, oldEmail)
		if err != nil {
			return errors.Wrap(err, "UserRepo.UpdateUser.ExecContext: deleting verification token")
		}
	}

	return nil
}
