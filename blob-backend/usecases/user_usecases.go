package usecases

import (
	"context"

	"github.com/joaoleau/blob/models"
	"github.com/joaoleau/blob/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type UserUseCase struct {
	repository *repository.UserRepo
}

func NewUserUseCase(repo *repository.UserRepo) *UserUseCase {
	return &UserUseCase{
		repository: repo,
	}
}

func (u *UserUseCase) GetUserById(ctx context.Context, id string) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUseCase.GetUserById")
	defer span.Finish()

	user, err := u.repository.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserUseCase) GetUserByUsername(ctx context.Context, username string) (*models.UserWithBlobs, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUseCase.GetUserByUsername")
	defer span.Finish()

	userWithBlobs, err := u.repository.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return userWithBlobs, nil
}

func (u *UserUseCase) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUseCase.GetUserByEmail")
	defer span.Finish()

	user, err := u.repository.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserUseCase) UpdateUser(ctx context.Context, email string, userData models.User) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUseCase.UpdateUser")
	defer span.Finish()

	user, err := u.GetUserByEmail(ctx, email)
	if err != nil {
		return errors.Wrap(err, "UserUseCase.UpdateUser.GetUserByEmail")
	}
	if user == nil {
		return errors.New("user not found")
	}

	err = u.repository.UpdateUser(ctx, user.ID, userData)
	if err != nil {
		return errors.Wrap(err, "UserUseCase.UpdateUser.UpdateUser")
	}

	return nil
}