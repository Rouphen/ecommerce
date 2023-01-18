package repository_test

import (
	"context"
	"regexp"
	"testing"

	"ecommerce/product-service/pkg/domain"
	productrepo "ecommerce/product-service/pkg/repository"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Test_Product_Create(t *testing.T) {
	t.Run("sucess_product_create", func(t *testing.T) {
		product := domain.Product{
			Name:  "pc",
			Stock: 15,
			Price: 150000,
		}

		productId := int64(3)
		mock, repo := NewMockProductRepository()

		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta("INSERT INTO `products` (`name`,`stock`,`price`) VALUES (?,?,?)")).
			WithArgs(product.Name, product.Stock, product.Price).
			WillReturnResult(sqlmock.NewResult(productId, 1))
		mock.ExpectCommit()

		err := repo.Create(context.TODO(), &product)

		assert.Nil(t, err)
		assert.Equal(t, productId, product.Id)
	})
}

func Test_Product_GetByID(t *testing.T) {
	t.Run("success_get_by_id", func(t *testing.T) {
		mock, repo := NewMockProductRepository()

		productId := int64(1)
		rows := sqlmock.NewRows([]string{"id", "name", "stock", "price"}).
			AddRow(productId, "pc", 15, 15000)

		mock.ExpectQuery("SELECT(.*)").
			WithArgs(productId).
			WillReturnRows(rows)

		product, err := repo.GetByID(context.TODO(), 1)

		assert.Nil(t, err)
		assert.Equal(t, productId, product.Id)
	})

	t.Run("failed_get_by_id", func(t *testing.T) {
		mock, repo := NewMockProductRepository()

		productId := int64(1)
		rows := sqlmock.NewRows([]string{"id", "name", "stock", "price"})

		mock.ExpectQuery("SELECT(.*)").
			WithArgs(productId).
			WillReturnRows(rows)

		_, err := repo.GetByID(context.TODO(), 1)

		assert.NotNil(t, err)
	})
}

func Test_Product_Save(t *testing.T) {
	t.Run("sucess_product_save", func(t *testing.T) {
		product := domain.Product{
			Name:  "pc",
			Stock: 15,
			Price: 150000,
		}

		productId := int64(3)
		mock, repo := NewMockProductRepository()

		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta("INSERT INTO `products` (`name`,`stock`,`price`) VALUES (?,?,?)")).
			WithArgs(product.Name, product.Stock, product.Price).
			WillReturnResult(sqlmock.NewResult(productId, 1))
		mock.ExpectCommit()

		err := repo.Save(context.TODO(), &product)

		assert.Nil(t, err)
		assert.Equal(t, productId, product.Id)
	})
}

func NewMockProductRepository() (sqlmock.Sqlmock, domain.ProductRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	DB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	return mock, productrepo.NewProdcutRepository(DB)
}
