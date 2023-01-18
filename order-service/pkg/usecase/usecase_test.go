package usecase_test

import (
	"context"
	domain "ecommerce/order-service/pkg/domain"
	"ecommerce/order-service/pkg/domain/mocks"
	ucase "ecommerce/order-service/pkg/usecase"

	"github.com/stretchr/testify/assert"

	"testing"

	"github.com/stretchr/testify/mock"
)

func Test_Create(t *testing.T) {
	repo, usecase := initalizeRepoAndUsecase()

	repo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Order")).Return(nil)

	mockOrder := &domain.Order{}
	response := usecase.Create(context.TODO(), mockOrder)

	assert.Equal(t, "", response.Error)
}

func Test_Delete(t *testing.T) {
	repo, usecase := initalizeRepoAndUsecase()

	t.Run("success_order_delete", func(t *testing.T) {
		repo.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(nil)

		response := usecase.Delete(context.TODO(), int64(0))

		assert.Equal(t, "", response.Error)
	})

	t.Run("falied_order_delete", func(t *testing.T) {
		repo.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return("not existed")

		response := usecase.Delete(context.TODO(), int64(0))

		assert.NotNil(t, response.Error)
	})
}

func initalizeRepoAndUsecase() (*mocks.OrderRepository, domain.OrderUsecase) {
	mockOrderRepo := new(mocks.OrderRepository)
	orderUcase := ucase.NewOrderUsecae(mockOrderRepo)

	return mockOrderRepo, orderUcase
}
