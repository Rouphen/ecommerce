package usecase

import (
	"context"
	"ecommerce/auth-service/pkg/domain"
	"ecommerce/auth-service/pkg/utils"
	"net/http"
)

type authUsecase struct {
	repo domain.UserRepository
	jwt  utils.JwtWrapper
}

func NewAuthUsecase(repo domain.UserRepository, jwt utils.JwtWrapper) domain.AuthUsecase {
	return &authUsecase{
		repo: repo,
		jwt:  jwt,
	}
}

func (a *authUsecase) Register(ctx context.Context, user *domain.User) domain.RegisterResponse {
	if _, err := a.repo.GetByEmail(ctx, user.Email); err == nil {
		return domain.RegisterResponse{
			Status: http.StatusConflict,
			Error:  "E-Mail already exists",
		}

	}

	user.Password = utils.HashPassword(user.Password)
	if err := a.repo.Create(ctx, user); err != nil {
		return domain.RegisterResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}
	}

	return domain.RegisterResponse{
		Status: http.StatusCreated,
	}
}

func (a *authUsecase) Login(ctx context.Context, req *domain.User) domain.LoginResponse {
	user, err := a.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return domain.LoginResponse{
			Status: http.StatusNotFound,
			Error:  "User not found",
		}
	}

	match := utils.CheckPasswordHash(req.Password, user.Password)
	if !match {
		return domain.LoginResponse{
			Status: http.StatusNotFound,
			Error:  "User not found",
		}
	}

	token, _ := a.jwt.GenerateToken(*user)
	// if err != nil {
	// 	return domain.LoginResponse{
	// 		Status: http.StatusNotFound,
	// 		Error:  "User not found",
	// 	}
	// }

	return domain.LoginResponse{
		Status: http.StatusOK,
		Token:  token,
	}
}

func (a *authUsecase) Validate(ctx context.Context, req *domain.User) domain.ValidateResponse {
	_, err := a.jwt.ValidateToken(req.Token)
	if err != nil {
		return domain.ValidateResponse{
			Status: http.StatusBadRequest,
		}
	}

	user, err := a.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return domain.ValidateResponse{
			Status: http.StatusNotFound,
			Error:  "User not found",
		}
	}

	return domain.ValidateResponse{
		Status: http.StatusOK,
		UserId: user.Id,
	}
}
