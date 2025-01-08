package test_for_analytic

import (
	"cachManagerApp/app/db/models"
	expenses "cachManagerApp/app/internal/methods-for-analytic/methods-for-expenses"
	"cachManagerApp/app/internal/tests/mocks"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

func setupMockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	// Оборачиваем sqlmock в GORM
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	return gormDB, mock, nil
}

func TestExpenseDayAnalytic(t *testing.T) {
	mocker := new(mocks.ExpenseAnalyticHandler)
	t.Run("AssertExpect", func(t *testing.T) {
		mocker.On("ExpenseDayAnalytic", mock.Anything).Return([]models.Transactions{}, nil)
		_, err := mocker.ExpenseDayAnalytic(tgbotapi.Update{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		mocker.AssertExpectations(t)
	})

	t.Run("TestExpDay", func(t *testing.T) {
		expectedTransactions := []models.Transactions{
			{
				ID:            1,
				TelegramID:    12345,
				OperationType: false,
				CreatedAt:     time.Now(),
				Quantities:    100,
			},
		}
		mocker.On("ExpenseDayAnalytic", mock.Anything).Return(expectedTransactions, nil)
		transactions, err := mocker.ExpenseDayAnalytic(Update)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		for i, tx := range transactions {
			if tx.ID != expectedTransactions[i].ID {
				t.Errorf("expected transaction ID %d, got %d", expectedTransactions[i].ID, tx.ID)
			}
			if tx.TelegramID != expectedTransactions[i].TelegramID {
				t.Errorf("expected TelegramID %d, got %d", expectedTransactions[i].TelegramID, tx.TelegramID)
			}
			if tx.OperationType != expectedTransactions[i].OperationType {
				t.Errorf("expected OperationType %v, got %v", expectedTransactions[i].OperationType, tx.OperationType)
			}
			if tx.Quantities != expectedTransactions[i].Quantities {
				t.Errorf("expected Amount %v, got %v", expectedTransactions[i].Quantities, tx.Quantities)
			}
		}
	})

}

func TestWeekAnalytic(t *testing.T) {
	mocker := new(mocks.ExpenseAnalyticHandler)
	result := make(map[string]uint64)
	t.Run("AssertExpect", func(t *testing.T) {
		mocker.On("ExpenseWeekAnalytic", mock.Anything).Return(result, nil)
		_, err := mocker.ExpenseWeekAnalytic(tgbotapi.Update{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		mocker.AssertExpectations(t)
	})

	t.Run("TestExpWeek", func(t *testing.T) {
		// Настраиваем мок-базу данных
		gormDB, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("failed to setup mock DB: %v", err)
		}

		expHandler := &expenses.ExpensesHandler{DB: gormDB}

		// Подготовим ожидаемые данные
		mockRows := sqlmock.NewRows([]string{"category", "value"}).
			AddRow("Food", 100).
			AddRow("Transport", 50)

		// Ожидаем SQL-запрос
		mock.ExpectQuery(`SELECT category, SUM.*FROM "transactions" WHERE .* GROUP BY "category"`).
			WithArgs(Update.Message.Chat.ID, false, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(mockRows)

		// Тестируем функцию
		result, err := expHandler.ExpenseWeekAnalytic(Update)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Проверяем результат
		expected := map[string]uint64{
			"Food":      100,
			"Transport": 50,
		}

		if len(result) != len(expected) {
			t.Fatalf("expected result length %d, got %d", len(expected), len(result))
		}

		for k, v := range expected {
			if result[k] != v {
				t.Errorf("expected %s = %d, got %d", k, v, result[k])
			}
		}

		// Проверяем ожидания
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}

func TestMonthAnalytic(t *testing.T) {
	t.Run("AssertExpect", func(t *testing.T) {
		mocker := new(mocks.ExpenseAnalyticHandler)
		result := make(map[string]uint64)
		mocker.On("ExpenseMonthAnalytic", mock.Anything).Return(result, nil)
		_, err := mocker.ExpenseMonthAnalytic(tgbotapi.Update{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		mocker.AssertExpectations(t)
	})

	t.Run("TestExpMonth", func(t *testing.T) {
		gormDB, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("failed to setup mock DB: %v", err)
		}
		expHandler := &expenses.ExpensesHandler{DB: gormDB}
		mockRows := sqlmock.NewRows([]string{"category", "value"}).
			AddRow("Food", 100).
			AddRow("Transport", 50)

		mock.ExpectQuery(`SELECT category, SUM.*FROM "transactions" WHERE.* GROUP BY "category"`).
			WithArgs(Update.Message.Chat.ID, false, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(mockRows)

		result, err := expHandler.ExpenseMonthAnalytic(Update)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := map[string]uint64{
			"Food":      100,
			"Transport": 50,
		}
		if len(result) != len(expected) {
			t.Fatalf("expected result length %d, got %d", len(expected), len(result))
		}
		for k, v := range expected {
			if result[k] != v {
				t.Errorf("expected %s = %d, got %v", k, v, result)
			}
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}
