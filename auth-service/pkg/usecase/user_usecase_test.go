package usecase

import (
	"context"
	"ecommerce/auth-service/pkg/domain"
	"ecommerce/auth-service/pkg/domain/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Get(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)

	mockUserRepo.On("Get", mock.Anything, mock.AnythingOfType("int64")).
		Return(&domain.User{}, nil)
	ucase := NewUserUsecase(mockUserRepo)

	_, err := ucase.Get(context.TODO(), 0)
	assert.Nil(t, err)
}
