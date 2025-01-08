package methods_for_transactions

import (
	"cachManagerApp/app/db/models"
	"cachManagerApp/database"
	"log/slog"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var categoryTranslations = map[string]string{
	// Доходы
	"salary":            "Заработная плата",
	"additional_income": "Побочный доход",
	"business_income":   "Доход от бизнеса",
	"investment_income": "Доход от инвестиций",
	"state_payments":    "Гос. выплаты",
	"property_sales":    "Продажа имущества",
	"other_income":      "Прочие доходы",
	// Расходы
	"basic_expense":      "Бытовые траты",
	"regular_expense":    "Регулярные платежи",
	"clothes":            "Одежда",
	"health":             "Здоровье",
	"leisure_education":  "Досуг и образование",
	"investment_expense": "Инвестиции",
	"other_expense":      "Прочие расходы",
}

var categoryIncomes = map[string]bool{
	// Доходы
	"salary":            true,
	"additional_income": true,
	"business_income":   true,
	"investment_income": true,
	"state_payments":    true,
	"property_sales":    true,
	"other_income":      true,
	// Расходы
	"basic_expense":      false,
	"regular_expense":    false,
	"clothes":            false,
	"health":             false,
	"leisure_education":  false,
	"investment_expense": false,
	"other_expense":      false,
}

type TransactionsMethod struct{}

func (transactions *TransactionsMethod) PostTransactionWithComment(update tgbotapi.Update, category string, amount int64, comment string, log *slog.Logger) error {
	transaction := models.Transactions{
		TelegramID:    uint64(update.Message.Chat.ID),
		CreatedAt:     time.Now(),
		OperationType: categoryIncomes[category], // здесь забираем тип приход или расход
		Quantities:    uint64(amount),
		Category:      categoryTranslations[category], // переводим здесь на русский
		Description:   comment,
	}

	if err := database.DB.Create(&transaction).Error; err != nil {
		log.Error("Ошибка добавления новой транзакции: %v", "err", err)
		return err
	}

	log.Info("Транзакция успешно добавлена", "transaction", transaction)
	return nil
}
