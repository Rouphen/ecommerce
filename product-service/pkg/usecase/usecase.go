package usecase

import (
	"context"
	"ecommerce/product-service/pkg/domain"
	"net/http"
)

type productUsecase struct {
	productRepo domain.ProductRepository
	stockRepo   domain.StockDecreaseLogsRepository
}

func NewProductUsecase(productRepo domain.ProductRepository,
	stockRepo domain.StockDecreaseLogsRepository) domain.ProductUsecase {

	return &productUsecase{
		productRepo: productRepo,
		stockRepo:   stockRepo,
	}
}

func (p *productUsecase) Create(ctx context.Context, product *domain.Product) domain.ProdcutResponse {
	err := p.productRepo.Create(ctx, product)

	var response domain.ProdcutResponse
	response.Status = http.StatusCreated
	if err != nil {
		response.Status = http.StatusConflict
		response.Error = err.Error()
	}

	return response
}

func (p *productUsecase) GetByID(ctx context.Context, reqId int64) (domain.Product, domain.ProdcutResponse) {
	var response domain.ProdcutResponse
	response.Status = http.StatusOK

	product, err := p.productRepo.GetByID(ctx, reqId)
	if err != nil {
		response.Status = http.StatusNotFound
		response.Error = err.Error()
	}

	return product, response
}

func (p *productUsecase) DecreaseStock(ctx context.Context, reqId, orderId int64) domain.ProdcutResponse {
	product, err := p.productRepo.GetByID(ctx, reqId)

	var response domain.ProdcutResponse
	response.Status = http.StatusOK
	if err != nil {
		response.Status = http.StatusNotFound
		response.Error = err.Error()

		return response
	}

	if product.Stock <= 0 {
		response.Status = http.StatusLengthRequired
		response.Error = "Stock too low"

		return response
	}

	stocklog, err := p.stockRepo.GetByOrderId(ctx, orderId)
	if err == nil {
		response.Status = http.StatusConflict
		response.Error = "Stock already decreased"
		return response
	}

	product.Stock = product.Stock - 1
	p.productRepo.Save(ctx, &product)

	stocklog.OrderId = orderId
	stocklog.ProductRefer = product.Id
	p.stockRepo.Create(ctx, &stocklog)

	response.Status = http.StatusOK

	return response
}
