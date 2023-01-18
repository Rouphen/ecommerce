package repository

import (
	"context"

	domain "ecommerce/product-service/pkg/domain"

	"gorm.io/gorm"
)

type productRepositoryImpl struct {
	db *gorm.DB
}

func NewProdcutRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepositoryImpl{
		db: db,
	}
}

func (p *productRepositoryImpl) Create(ctx context.Context, product *domain.Product) error {
	result := p.db.Create(product)

	return result.Error
}

func (p *productRepositoryImpl) GetByID(ctx context.Context, reqId int64) (domain.Product, error) {
	var product domain.Product
	result := p.db.First(&product, reqId)
	if result.Error != nil {
		return product, result.Error
	}

	return product, nil
}

func (p *productRepositoryImpl) Save(ctx context.Context, product *domain.Product) error {
	result := p.db.Save(product)
	return result.Error
}
