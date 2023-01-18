package usecase

import (
	"context"
	"ecommerce/auth-service/pkg/domain"
	"ecommerce/auth-service/pkg/domain/mocks"
	"ecommerce/auth-service/pkg/utils"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
)

func Test_UserLogin(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockUser := &domain.User{
		Id:       1,
		Email:    "test@12.com",
		Password: "password1",
	}

	mockHashUser := &domain.User{
		Id:       1,
		Email:    "test@12.com",
		Password: "password1",
	}

	mockTokenUser := &domain.User{
		Id:       1,
		Email:    "",
		Password: "password1",
	}
	mockHashUser.Password = utils.HashPassword(mockHashUser.Password)
	mockTokenUser.Password = utils.HashPassword(mockTokenUser.Password)

	cases := []struct {
		Name       string
		Expect     int64
		User       *domain.User
		GetByEmail *mock.Call
	}{
		{"User Not found", http.StatusNotFound, mockUser,
			mockUserRepo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).
				Return(&domain.User{}, errors.New("Not found"))},

		{"Password error", http.StatusNotFound, mockUser,
			mockUserRepo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).
				Return(&domain.User{}, nil)},

		{"Token invalid", http.StatusNotFound, mockUser,
			mockUserRepo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).
				Return(mockTokenUser, nil)},

		{"User is  valid", http.StatusOK, mockUser,
			mockUserRepo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).
				Return(mockHashUser, nil)},
	}

	jwt := utils.JwtWrapper{
		SecretKey:       "r43t18sc",
		Issuer:          "auth-service",
		ExpirationHours: 24 * 365,
	}

	ucase := NewAuthUsecase(mockUserRepo, jwt)

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			tt.GetByEmail.Once()
			response := ucase.Login(context.TODO(), tt.User)
			assert.Equal(t, int64(tt.Expect), response.Status)
		})
	}
}

func Test_UserRegister(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockUser := &domain.User{
		Id:       1,
		Email:    "test@12.com",
		Password: "password1",
	}

	cases := []struct {
		Name       string
		Expect     int64
		User       *domain.User
		GetByEmail *mock.Call
		Create     *mock.Call
	}{
		{"Password is OK", http.StatusCreated, mockUser,
			mockUserRepo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).
				Return(&domain.User{}, errors.New("Not found")),
			mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).
				Return(nil),
		},

		{"Password error", http.StatusConflict, mockUser,
			mockUserRepo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).
				Return(mockUser, nil),
			mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).
				Return(errors.New("failed to create")),
		},

		{"User Not found", http.StatusBadRequest, mockUser,
			mockUserRepo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).
				Return(&domain.User{}, errors.New("Not found")),
			mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).
				Return(errors.New("Not found")),
		},
	}

	jwt := utils.JwtWrapper{
		SecretKey:       "r43t18sc",
		Issuer:          "auth-service",
		ExpirationHours: 24 * 365,
	}

	ucase := NewAuthUsecase(mockUserRepo, jwt)

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			tt.GetByEmail.Once()
			tt.Create.Once()
			response := ucase.Register(context.TODO(), tt.User)
			assert.Equal(t, int64(tt.Expect), response.Status)
		})
	}
}

func Test_UserValidate(t *testing.T) {
	jwt := utils.JwtWrapper{
		SecretKey:       "r43t18sc",
		Issuer:          "auth-service",
		ExpirationHours: 24 * 365,
	}

	jwt_expirated := utils.JwtWrapper{
		SecretKey:       "r43t18s",
		Issuer:          "auth-service",
		ExpirationHours: 0,
	}
	mockTokenUser := &domain.User{
		Id: 1,
	}
	tokenMock, _ := jwt_expirated.GenerateToken(*mockTokenUser)
	mockTokenUser.Token = tokenMock

	mockUserRepo := new(mocks.UserRepository)
	mockUser := &domain.User{
		Id:       1,
		Email:    "test@12.com",
		Password: "password1",
	}
	token, _ := jwt.GenerateToken(*mockUser)
	mockUser.Token = token

	cases := []struct {
		Name       string
		Expect     int64
		User       *domain.User
		GetByEmail *mock.Call
	}{
		{"Token is OK", http.StatusOK, mockUser,
			mockUserRepo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).
				Return(mockUser, nil),
		},

		{"Token is warn", http.StatusBadRequest, mockTokenUser,
			mockUserRepo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).
				Return(&domain.User{}, errors.New("Not found")),
		},

		{"User Not found", http.StatusNotFound, mockUser,
			mockUserRepo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).
				Return(mockUser, errors.New("Not found")),
		},
	}

	ucase := NewAuthUsecase(mockUserRepo, jwt)

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			tt.GetByEmail.Once()
			response := ucase.Validate(context.TODO(), tt.User)
			assert.Equal(t, int64(tt.Expect), response.Status)
		})
	}
}
