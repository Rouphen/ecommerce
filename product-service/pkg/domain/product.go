package domain

import "context"

type Product struct {
	Id                int64            `json:"id" gorm:"primaryKey"`
	Name              string           `json:"name"`
	Stock             int64            `json:"stock"`
	Price             int64            `json:"price"`
	stockDecreaseLogs StockDecreaseLog `gorm:"foreignKey:ProductRefer"`
}

type ProdcutResponse struct {
	Status int64  `json:"status"`
	Error  string `json:"error"`
}

type ProductUsecase interface {
	Create(ctx context.Context, product *Product) ProdcutResponse
	GetByID(ctx context.Context, eqId int64) (Product, ProdcutResponse)
	DecreaseStock(ctx context.Context, reqId, orderId int64) ProdcutResponse
}

type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	GetByID(ctx context.Context, reqId int64) (Product, error)
	Save(ctx context.Context, product *Product) error
}
