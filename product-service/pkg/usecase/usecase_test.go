package usecase_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"ecommerce/product-service/pkg/domain"
	mocks "ecommerce/product-service/pkg/domain/mocks"
	ucase "ecommerce/product-service/pkg/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Create(t *testing.T) {
	mockProduct := domain.Product{
		Name:  "test_product_name",
		Stock: 2,
		Price: 5,
	}

	t.Run("success", func(t *testing.T) {
		mockProductRepo, _, productUsecase := initalizeRepoAndUsecase()
		mockProductRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Product")).Return(nil).Once()
		response := productUsecase.Create(context.TODO(), &mockProduct)

		assert.Equal(t, response.Status, int64(http.StatusCreated))
	})

	t.Run("failed", func(t *testing.T) {
		mockProductRepo, _, productUsecase := initalizeRepoAndUsecase()
		mockProductRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Product")).Return(errors.New("conflict")).Once()
		response := productUsecase.Create(context.TODO(), &mockProduct)

		assert.Equal(t, response.Status, int64(http.StatusConflict))
	})
}

func Test_GetByID(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		mockProductRepo, _, productUsecase := initalizeRepoAndUsecase()

		mockProductRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(domain.Product{}, nil).Once()
		_, response := productUsecase.GetByID(context.TODO(), int64(1))

		assert.Equal(t, response.Status, int64(http.StatusOK))
	})

	t.Run("failed", func(t *testing.T) {
		mockProductRepo, _, productUsecase := initalizeRepoAndUsecase()

		mockProductRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(domain.Product{}, errors.New("not existed"))

		_, response := productUsecase.GetByID(context.TODO(), 1)

		assert.Equal(t, response.Status, int64(http.StatusNotFound))
	})

}

func Test_DecreaseStock(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockProduct := domain.Product{
			Name:  "test_product_name",
			Stock: 2,
			Price: 5,
		}

		mockProductRepo, mockStockLogRepo, productUsecase := initalizeRepoAndUsecase()
		mockProductRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockProduct, nil)
		mockStockLogRepo.On("GetByOrderId", mock.Anything, mock.AnythingOfType("int64")).Return(domain.StockDecreaseLog{}, errors.New("Not found"))
		mockProductRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.Product")).Return(nil)
		mockStockLogRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.StockDecreaseLog")).Return(nil)

		response := productUsecase.DecreaseStock(context.TODO(), 1, 1)

		assert.Equal(t, response.Status, int64(http.StatusOK))
	})

	t.Run("failed--no product", func(t *testing.T) {
		mockProductRepo, _, productUsecase := initalizeRepoAndUsecase()
		mockProductRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(domain.Product{}, errors.New("not found"))

		response := productUsecase.DecreaseStock(context.TODO(), 1, 1)

		assert.Equal(t, response.Status, int64(http.StatusNotFound))
	})

	t.Run("failed--size invalid", func(t *testing.T) {
		mockProductRepo, _, productUsecase := initalizeRepoAndUsecase()
		mockProductRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(domain.Product{}, nil)

		response := productUsecase.DecreaseStock(context.TODO(), 1, 1)

		assert.Equal(t, response.Status, int64(http.StatusLengthRequired))
	})

	t.Run("failed--orderid is existed", func(t *testing.T) {
		mockProduct := domain.Product{
			Name:  "test_product_name",
			Stock: 2,
			Price: 5,
		}

		mockProductRepo, mockStockLogRepo, productUsecase := initalizeRepoAndUsecase()
		mockProductRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockProduct, nil)
		mockStockLogRepo.On("GetByOrderId", mock.Anything, mock.AnythingOfType("int64")).Return(domain.StockDecreaseLog{}, nil)

		response := productUsecase.DecreaseStock(context.TODO(), 1, 1)

		assert.Equal(t, response.Status, int64(http.StatusConflict))
	})
}

func initalizeRepoAndUsecase() (*mocks.ProductRepository, *mocks.StockDecreaseLogsRepository, domain.ProductUsecase) {
	mockProductRepo := new(mocks.ProductRepository)
	mockStockLogRepo := new(mocks.StockDecreaseLogsRepository)
	productUsecase := ucase.NewProductUsecase(mockProductRepo, mockStockLogRepo)

	return mockProductRepo, mockStockLogRepo, productUsecase
}
