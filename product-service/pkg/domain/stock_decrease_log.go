package domain

import "context"

type StockDecreaseLog struct {
	Id           int64 `json:"id" gorm:"primaryKey"`
	OrderId      int64 `json:"order_id"`
	ProductRefer int64 `json:"product_id"`
}

type StockDecreaseLogsRepository interface {
	Create(ctx context.Context, log *StockDecreaseLog) error
	GetByOrderId(ctx context.Context, orderId int64) (StockDecreaseLog, error)
}
