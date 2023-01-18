package usecase

import (
	"context"
	domain "ecommerce/order-service/pkg/domain"
	"net/http"
)

type orderUsecaseImpl struct {
	repo domain.OrderRepository
}

func NewOrderUsecae(repo domain.OrderRepository) domain.OrderUsecase {
	return &orderUsecaseImpl{
		repo: repo,
	}
}

func (o *orderUsecaseImpl) Create(ctx context.Context, order *domain.Order) *domain.OrderResponse {
	var response domain.OrderResponse
	response.Status = http.StatusOK

	if err := o.repo.Create(ctx, order); err != nil {
		response.Error = err.Error()
	}

	return &response
}

func (o *orderUsecaseImpl) Delete(ctx context.Context, orderID int64) *domain.OrderResponse {
	var response domain.OrderResponse
	if err := o.repo.Delete(ctx, orderID); err != nil {
		response.Error = err.Error()
	}

	return &response
}
