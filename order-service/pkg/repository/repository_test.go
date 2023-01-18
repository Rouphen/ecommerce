package repository_test

import (
	"context"
	"ecommerce/order-service/pkg/domain"
	orderrepo "ecommerce/order-service/pkg/repository"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Test_Create(t *testing.T) {
	mock, repo := NewMockProductRepository()
	mockOrder := &domain.Order{
		Price:     15000,
		ProductId: 2,
		UserId:    1,
	}

	orderId := int64(3)

	mock.ExpectBegin()
	mock.ExpectExec(
		regexp.QuoteMeta("INSERT INTO `orders` (`price`,`product_id`,`user_id`) VALUES (?,?,?)")).
		WithArgs(mockOrder.Price, mockOrder.ProductId, mockOrder.UserId).
		WillReturnResult(sqlmock.NewResult(orderId, 1))
	mock.ExpectCommit()

	err := repo.Create(context.TODO(), mockOrder)

	assert.Nil(t, err)
	assert.Equal(t, orderId, mockOrder.Id)
}

func Test_Delete(t *testing.T) {
	mock, repo := NewMockProductRepository()
	orderId := int64(4)

	mock.ExpectBegin()
	mock.ExpectExec(
		regexp.QuoteMeta("DELETE FROM `orders` WHERE order_id=?")).
		WithArgs(orderId).
		WillReturnResult(sqlmock.NewResult(orderId, 1))
	mock.ExpectCommit()

	err := repo.Delete(context.TODO(), orderId)
	assert.Nil(t, err)
}

func NewMockProductRepository() (sqlmock.Sqlmock, domain.OrderRepository) {
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

	return mock, orderrepo.NewOrderRepository(DB)
}
