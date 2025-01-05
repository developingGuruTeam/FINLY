package TgBot

import (
	"cachManagerApp/app/internal/methodsForTransaction"
	"cachManagerApp/app/internal/methodsForUser"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
)

// Обработка транзакции
func handleTransactionAction(bot *tgbotapi.BotAPI, update tgbotapi.Update, transResp TransactionResponse, log *slog.Logger) {
	chatID := update.Message.Chat.ID
	switch transResp.Action {
	// incomes
	case "salary":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Заработная плата"
		if err := transaction.PostIncome(update, category, log); err != nil {
			log.Info("Failed to save salary: %v", log.With("error", err))
		}
		msg := tgbotapi.NewMessage(chatID, "Заработная плата сохранена.")
		if _, err := bot.Send(msg); err != nil {
			log.Info("Failed to send salary message: %v", log.With("error", err))
		}

	case "additional_income":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Побочный доход"
		if err := transaction.PostIncome(update, category, log); err != nil {
			log.Error("Failed to save additional income: %v", log.With("error", err))
		}
		msg := tgbotapi.NewMessage(chatID, "Побочный доход сохранен.")
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send additional income message: %v", log.With("error", err))
		}

	case "business_income":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Доход от бизнеса"
		if err := transaction.PostIncome(update, category, log); err != nil {
			log.Error("Failed to save business income: %v", log.With("error", err))
		}
		msg := tgbotapi.NewMessage(chatID, "Доход от бизнеса сохранен.")
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send business income message: %v", log.With("error", err))
		}

	case "investment_income":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Доход от инвестиций"
		if err := transaction.PostIncome(update, category, log); err != nil {
			log.Error("Failed to save investment income: %v", log.With("error", err))
		}
		msg := tgbotapi.NewMessage(chatID, "Доход от инвестиций сохранен.")
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send investment income message: %v", log.With("error", err))
		}

	case "state_payments":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Гос. выплаты"
		if err := transaction.PostIncome(update, category, log); err != nil {
			log.Error("Failed to save investment income: %v", log.With("error", err))
		}
		msg := tgbotapi.NewMessage(chatID, "Доход от государства сохранен.")
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send state income message: %v", log.With("error", err))
		}

	case "property_sales":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Продажа имущества"
		if err := transaction.PostIncome(update, category, log); err != nil {
			log.Error("Failed to save investment income: %v", log.With("error", err))
		}
		msg := tgbotapi.NewMessage(chatID, "Доход от продажи имущества сохранен.")
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send investment income message: %v", log.With("error", err))
		}

	case "other_income":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Прочие доходы"
		if err := transaction.PostIncome(update, category, log); err != nil {
			log.Error("Failed to save other income: %v", log.With("error", err))
		}
		msg := tgbotapi.NewMessage(chatID, "Прочие доходы сохранены.")
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send other income message: %v", log.With("error", err))
		}
	// expenses
	case "basic_expense":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Бытовые траты"
		if err := transaction.PostExpense(update, category, log); err != nil {
			log.Error("Failed to save basic expense: %v", log.With("error", err))
		}
		msg := tgbotapi.NewMessage(chatID, "Сумма базовых трат сохранена.")
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send basic expense message: ", log.With("error", err))
		}

	case "regular_expense":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Регулярные платежи"
		if err := transaction.PostExpense(update, category, log); err != nil {
			log.Error("Failed to save regular expense:", log.With("error", err))
		}
		msg := tgbotapi.NewMessage(chatID, "Регулярный платеж сохранен.")
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send regular expense message:", log.With("error", err))
		}

	case "clothes":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Одежда"
		if err := transaction.PostExpense(update, category, log); err != nil {
			log.Error("Failed to save clothes", log.With("error", err))
		}
		msg := tgbotapi.NewMessage(chatID, "Расход на обновление гардероба сохранен.")
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send clothes message", log.With("error", err))
		}

	case "health":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Здоровье"
		if err := transaction.PostExpense(update, category, log); err != nil {
			log.Error("Failed to save health", log.With("error", err))
		}
		msg := tgbotapi.NewMessage(chatID, "Расход на поддержание здоровья сохранен.")
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send health message", log.With("error", err))
		}

	case "leisure_education":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Досуг и образование"
		if err := transaction.PostExpense(update, category, log); err != nil {
			log.Error("Failed to save leisure_education expense", log.With("error", err))
		}
		msg := tgbotapi.NewMessage(chatID, "Расход сохранен.")
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send leisure_education message", log.With("error", err))
		}

	case "investment_expense":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Инвестиции"
		if err := transaction.PostExpense(update, category, log); err != nil {
			log.Error("Failed to save investment expense", log.With("error", err))
		}
		msg := tgbotapi.NewMessage(chatID, "Инвестиционный расход сохранен.")
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send investment expense message", log.With("error", err))
		}

	case "other_expense":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Прочие расходы"
		if err := transaction.PostExpense(update, category, log); err != nil {
			log.Error("Failed to save other expense", log.With("error", err))
		}
		msg := tgbotapi.NewMessage(chatID, "Прочие расходы сохранены.")
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send other expense message", log.With("error", err))
		}
	}
	mu.Lock()
	delete(transactionStates, chatID) // удаляем состояние после обработки
	mu.Unlock()

}

func handleUserAction(bot *tgbotapi.BotAPI, update tgbotapi.Update, userResp UserResponse, log *slog.Logger) {
	chatID := update.Message.Chat.ID

	switch userResp.Action {
	case "expense":
		amount := update.Message.Text
		msg := tgbotapi.NewMessage(chatID, "Сумма расхода "+amount+" сохранена.")
		if _, err := bot.Send(msg); err != nil {
			log.Error("Ошибка отправки сообщения о расходе", log.With("error", err))
		}

	case "income":
		amount := update.Message.Text
		msg := tgbotapi.NewMessage(chatID, "Сумма прихода "+amount+" сохранена.")
		if _, err := bot.Send(msg); err != nil {
			log.Error("Ошибка отправки сообщения о приходе", log.With("error", err))
		}

	case "change_name":
		// Обновление имени пользователя в БД
		user := methodsForUser.UserMethod{}
		if err := user.UpdateUserName(update); err != nil {
			log.Error("Ошибка обновления имени пользователя", log.With("error", err))
		}
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s! Ваше имя успешно изменено.", update.Message.Text))
		if _, err := bot.Send(msg); err != nil {
			log.Error("Ошибка отправки сообщения об изменении имени", log.With("error", err))
		}

	case "change_currency":
		user := methodsForUser.UserMethod{}
		if err := user.UpdateUserCurrency(update); err != nil {
			log.Error("Ошибка обновления валюты", log.With("error", err))
		}
		msg := tgbotapi.NewMessage(chatID, "Ваша валюта изменена.")
		if _, err := bot.Send(msg); err != nil {
			log.Error("Ошибка отправки сообщения об изменении валюты", log.With("error", err))
		}
	}

	mu.Lock()
	delete(userStates, chatID) // удаляем состояние после обработки
	mu.Unlock()
}
