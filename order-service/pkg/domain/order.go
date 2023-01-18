package domain

import "context"

type Order struct {
	Id        int64 `json:"id" gorm:"primaryKey"`
	Price     int64 `json:"price"`
	ProductId int64 `json:"product_id"`
	UserId    int64 `json:"user_id"`
}

type OrderResponse struct {
	Status int64
	Error  string
}

type OrderUsecase interface {
	Create(ctx context.Context, order *Order) *OrderResponse
	Delete(ctx context.Context, orderID int64) *OrderResponse
}

type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	Delete(ctx context.Context, orderId int64) error
}
