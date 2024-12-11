package methodsForTransaction

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TransactionsMethod struct{}

type TransactionsHandlers interface {
	PostIncome(update tgbotapi.Update, category string) error
	PostExpense(update tgbotapi.Update, category string) error
}
