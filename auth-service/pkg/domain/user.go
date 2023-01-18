package domain

import "context"

type User struct {
	Id       int64  `json:"id" gorm:"primaryKey"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

// type UserLogin struct {
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }

// type UserRegister struct {
// 	Email     string `json:"email"`
// 	Password  string `json:"password"`
// 	FirstName string `json:"first_name"`
// 	LastName  string `json:"last_name"`
// }

// type UserToken struct {
// 	Token string `json:"token"`
// }

type UserUsecase interface {
	Get(ctx context.Context, id int64) (*User, error)
	GetAll(ctx context.Context) ([]*User, error)
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type UserRepository interface {
	Get(ctx context.Context, id int64) (*User, error)
	GetAll(ctx context.Context) ([]*User, error)
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
}
