package usecase

import (
	"context"

	"ecommerce/auth-service/pkg/domain"
)

type userUsecase struct {
	repo domain.UserRepository
}

func NewUserUsecase(repo domain.UserRepository) domain.UserUsecase {
	return &userUsecase{
		repo: repo,
	}
}

func (u *userUsecase) Get(ctx context.Context, id int64) (*domain.User, error) {
	return u.repo.Get(ctx, id)
}

func (u *userUsecase) GetAll(ctx context.Context) ([]*domain.User, error) {
	return u.repo.GetAll(ctx)
}

func (u *userUsecase) Create(ctx context.Context, user *domain.User) error {
	return u.repo.Create(ctx, user)
}

func (u *userUsecase) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return u.repo.GetByEmail(ctx, email)
}
