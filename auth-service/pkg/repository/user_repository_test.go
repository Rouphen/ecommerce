package repository

import (
	"context"
	domain "ecommerce/auth-service/pkg/domain"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Test_User_Create(t *testing.T) {
	t.Run("sucess_user_create", func(t *testing.T) {
		mockUser := &domain.User{
			Id:       0,
			Email:    "xxx@yyy.com",
			Password: "password1",
		}

		userId := int64(3)
		mock, repo := NewMockUserRepository()

		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta("INSERT INTO `users` (`email`,`password`,`token`) VALUES (?,?,?)")).
			WithArgs(mockUser.Email, mockUser.Password, mockUser.Token).
			WillReturnResult(sqlmock.NewResult(userId, 1))
		mock.ExpectCommit()

		err := repo.Create(context.TODO(), mockUser)

		assert.Nil(t, err)
		assert.Equal(t, userId, mockUser.Id)
	})
}

func Test_GetAll(t *testing.T) {
	t.Run("sucess_get_all", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "password", "token"}).
			AddRow(0, "xxx0@yyy.com", "password1", "").
			AddRow(1, "xxx1@yyy.com", "password2", "")

		mock, repo := NewMockUserRepository()

		mock.ExpectQuery(
			regexp.QuoteMeta("SELECT * FROM `users`")).
			WillReturnRows(rows)

		users, err := repo.GetAll(context.TODO())

		assert.Nil(t, err)
		assert.Equal(t, 2, len(users))
	})
}

func NewMockUserRepository() (sqlmock.Sqlmock, domain.UserRepository) {
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

	return mock, NewUserRepository(DB)
}
