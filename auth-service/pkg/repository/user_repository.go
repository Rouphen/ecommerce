package repository

import (
	"context"
	domain "ecommerce/auth-service/pkg/domain"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) Get(ctx context.Context, id int64) (*domain.User, error) {
	var user domain.User
	if result := u.db.Where(&domain.User{Id: id}).First(&user); result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (u *userRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	var users []*domain.User
	if result := u.db.Find(&users); result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func (u *userRepository) Create(ctx context.Context, user *domain.User) error {
	result := u.db.Create(user)
	return result.Error
}

func (u *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if result := u.db.Where(&domain.User{Email: email}).First(&user); result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
