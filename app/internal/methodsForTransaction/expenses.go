package methodsForTransaction

import (
	"cachManagerApp/app/db/models"
	"log/slog"

	"cachManagerApp/database"
	"errors"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (transactions *TransactionsMethod) PostExpense(update tgbotapi.Update, category string, log *slog.Logger) error {
	var sum int
	var err error
	var userText string

	sum, err = strconv.Atoi(update.Message.Text)
	if err != nil || sum <= 0 {
		log.Info("Ошибка преобразования суммы: %v", "err", err)
		return errors.New("введено некорректное число")
	}
	userText = "" // Описание всегда пустое

	transaction := models.Transactions{
		TelegramID:    uint64(update.Message.Chat.ID),
		CreatedAt:     time.Now(),
		OperationType: false,
		Quantities:    uint64(sum),
		Category:      category,
		Description:   userText,
	}

	var transactionExist models.Transactions
	res := database.DB.Where("telegram_id = ? AND created_at = ?", transaction.TelegramID, transaction.CreatedAt).First(&transactionExist).Error
	if res == nil {
		log.Debug("Транзакция существует")
		return errors.New("transaction already exists")
	}

	if err := database.DB.Create(&transaction).Error; err != nil {
		log.Error("Ошибка добавления новой транзакции: %v", "err", err)
		return err

	}
	log.Info("Транзакция успешно добавлена", "transaction", transaction)
	return nil
}
