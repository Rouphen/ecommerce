package repository

import (
	"context"
	"ecommerce/order-service/pkg/domain"

	"gorm.io/gorm"
)

type orderRepositoryImpl struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) domain.OrderRepository {
	return &orderRepositoryImpl{
		db: db,
	}
}

func (o *orderRepositoryImpl) Create(ctx context.Context, order *domain.Order) error {
	result := o.db.Create(order)
	return result.Error
}

func (o *orderRepositoryImpl) Delete(ctx context.Context, orderId int64) error {
	result := o.db.Where("order_id=?", orderId).Delete(&domain.Order{})

	return result.Error
}
