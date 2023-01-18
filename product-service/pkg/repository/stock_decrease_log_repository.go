package repository

import (
	"context"

	domain "ecommerce/product-service/pkg/domain"

	"gorm.io/gorm"
)

type stockDecreaseLogsRepositoryImpl struct {
	db *gorm.DB
}

func NewStockDecreaseLogsRepository(db *gorm.DB) domain.StockDecreaseLogsRepository {
	return &stockDecreaseLogsRepositoryImpl{
		db: db,
	}
}

func (s *stockDecreaseLogsRepositoryImpl) Create(ctx context.Context, log *domain.StockDecreaseLog) error {
	result := s.db.Create(log)
	return result.Error
}

func (s *stockDecreaseLogsRepositoryImpl) GetByOrderId(ctx context.Context, orderId int64) (domain.StockDecreaseLog, error) {
	var log domain.StockDecreaseLog
	result := s.db.Where(&domain.StockDecreaseLog{OrderId: orderId}).First(&log)
	if result.Error != nil {
		return log, result.Error
	}

	return log, nil
}
