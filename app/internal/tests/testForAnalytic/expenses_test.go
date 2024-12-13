package testForAnalytic

import (
	"cachManagerApp/app/db/models"
	"cachManagerApp/app/internal/tests/mocks"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
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
		if len(transactions) != len(expectedTransactions) {
			t.Errorf("expected %d transactions, got %d", len(expectedTransactions), len(transactions))
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
