package test_for_user_methods

import (
	"cachManagerApp/app/db/models"
	"cachManagerApp/app/internal/tests/mocks"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Update = tgbotapi.Update{
	Message: &tgbotapi.Message{
		Chat: &tgbotapi.Chat{
			ID: 12345,
		},
		From: &tgbotapi.User{
			UserName: "TestUser",
		},
	},
}

func TestUsersHandlers(t *testing.T) {
	mocker := new(mocks.UsersHandlers)
	// тест функции PostUser
	t.Run("TestPostUser", func(t *testing.T) {
		mocker.On("PostUser", Update).Return(nil)
		err := mocker.PostUser(Update)
		assert.NoError(t, err)
		mocker.AssertCalled(t, "PostUser", Update)
	},
	)

	// тест функции TestUpdateUserName
	t.Run("TestUpdateUserName", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		gormDB, err := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "users" SET "name"=\$1 WHERE telegram_id = \$2`).
			WithArgs("NewUserName", 12345).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		result := gormDB.Model(&models.Users{}).
			Where("telegram_id = ?", 12345).
			Update("name", "NewUserName")
		assert.NoError(t, result.Error)
		assert.Equal(t, result.RowsAffected, int64(1))
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Не все ожидания были выполнены: %s", err)
		}
	})
	// тест функции UpdateUserCurrency
	t.Run("UpdateUserCurrency", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		gormDB, err := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "users" SET "currency"=\$1 WHERE telegram_id = \$2`).
			WithArgs("USD", 12345).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		result := gormDB.Model(&models.Users{}).
			Where("telegram_id = ?", 12345).
			Update("currency", "USD")

		assert.NoError(t, result.Error)
		assert.Equal(t, result.RowsAffected, int64(1))
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
		}
	})
}
