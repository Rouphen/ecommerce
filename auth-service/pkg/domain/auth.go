package domain

import (
	"context"
)

type AuthResponse struct {
	Id     int64  `json:"id"`
	Status int64  `json:"status"`
	Error  string `json:"error"`
	Token  string `json:"token"`
}

type RegisterResponse struct {
	Status int64
	Error  string
}

type LoginResponse struct {
	Status int64
	Error  string
	Token  string
}

type ValidateResponse struct {
	Status int64
	Error  string
	UserId int64
}

type AuthUsecase interface {
	Register(ctx context.Context, user *User) RegisterResponse
	Login(ctx context.Context, user *User) LoginResponse
	Validate(ctx context.Context, user *User) ValidateResponse
}
