package usecases

import (
	"context"
	"github.com/joaoleau/blob/repository"
	"github.com/joaoleau/blob/models"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
)

type AuthUseCase struct {
	repository repository.AuthRepo
}

func NewAuthUseCase(repo repository.AuthRepo) AuthUseCase {
	return AuthUseCase{
		repository: repo,
	}
}

// Create new user
func (u *AuthUseCase) Register(ctx context.Context, user *models.User) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AuthUseCase.Register")
	defer span.Finish()

	existsUser, err := u.repository.FindByEmail(ctx, user)
	if existsUser != nil || err == nil {
		return nil, err
	}

	if err = user.PrepareCreate(); err != nil {
		return nil, err
	}

	createdUser, err := u.repository.Register(ctx, user)
	if err != nil {
		return nil, err
	}
	createdUser.SanitizePassword()

	return createdUser, nil
}


// Update existing user
func (u *AuthUseCase) Update(ctx context.Context, user *models.User) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AuthUseCase.Update")
	defer span.Finish()

	if err := user.PrepareUpdate(); err != nil {
		return nil, err
	}

	updatedUser, err := u.repository.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	updatedUser.SanitizePassword()

	return updatedUser, nil
}


// Delete new user
func (u *AuthUseCase) Delete(ctx context.Context, userID uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AuthUseCase.Delete")
	defer span.Finish()

	if err := u.repository.Delete(ctx, userID); err != nil {
		return err
	}

	return nil
}


// Get user by id
func (u *AuthUseCase) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AuthUseCase.GetByID")
	defer span.Finish()

	user, err := u.repository.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.SanitizePassword()

	return user, nil
}
