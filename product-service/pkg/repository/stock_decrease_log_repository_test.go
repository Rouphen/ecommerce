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

func Test_StockDecreaseLog_Create(t *testing.T) {
	t.Run("sucess_StockDecreaseLog_create", func(t *testing.T) {
		stockDecreaseLog := domain.StockDecreaseLog{
			OrderId:      1,
			ProductRefer: 1,
		}

		DecreaseId := int64(3)
		mock, repo := NewMockSockDecreaseLogsRepository()

		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta("INSERT INTO `stock_decrease_logs` (`order_id`,`product_refer`) VALUES (?,?)")).
			WithArgs(stockDecreaseLog.OrderId, stockDecreaseLog.ProductRefer).
			WillReturnResult(sqlmock.NewResult(DecreaseId, 1))
		mock.ExpectCommit()

		err := repo.Create(context.TODO(), &stockDecreaseLog)

		assert.Nil(t, err)
		assert.Equal(t, DecreaseId, stockDecreaseLog.Id)
	})
}

func Test_StockDecreaseLog_GetByOrderId(t *testing.T) {
	t.Run("sucess_StockDecreaseLog_GetByOrderId", func(t *testing.T) {
		mock, repo := NewMockSockDecreaseLogsRepository()

		orderId := int64(1)
		rows := sqlmock.NewRows([]string{"id", "order_id", "product_refer"}).
			AddRow(0, orderId, 1)

		mock.ExpectQuery(
			regexp.QuoteMeta("SELECT * FROM `stock_decrease_logs` WHERE `stock_decrease_logs`.`order_id` = ?")).
			WithArgs(orderId).
			WillReturnRows(rows)

		stockDescreaseLog, err := repo.GetByOrderId(context.TODO(), orderId)

		assert.Nil(t, err)
		assert.Equal(t, orderId, stockDescreaseLog.OrderId)
	})

	t.Run("failed_StockDecreaseLog_GetByOrderId", func(t *testing.T) {
		mock, repo := NewMockSockDecreaseLogsRepository()

		orderId := int64(1)
		rows := sqlmock.NewRows([]string{"id", "order_id", "product_refer"})

		mock.ExpectQuery(
			regexp.QuoteMeta("SELECT * FROM `stock_decrease_logs` WHERE `stock_decrease_logs`.`order_id` = ?")).
			WithArgs(orderId).
			WillReturnRows(rows)

		_, err := repo.GetByOrderId(context.TODO(), orderId)

		assert.NotNil(t, err)
	})
}

func NewMockSockDecreaseLogsRepository() (sqlmock.Sqlmock, domain.StockDecreaseLogsRepository) {
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

	return mock, productrepo.NewStockDecreaseLogsRepository(DB)
}
