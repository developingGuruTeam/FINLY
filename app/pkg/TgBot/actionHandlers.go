package TgBot

import (
	"cachManagerApp/app/internal/methodsForTransaction"
	"cachManagerApp/app/internal/methodsForUser"
	"fmt"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
			log.Error("Failed to save salary:", slog.Any("error", err))
		}
		msg := tgbotapi.NewMessage(chatID, "Заработная плата сохранена.")
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send salary message:", slog.Any("error", err))
		}

	case "additional_income":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Побочный доход"
		if err := transaction.PostIncome(update, category, log); err != nil {
			log.Error("Failed to save additional income:", slog.Any("error", err))
		}
		msg := tgbotapi.NewMessage(chatID, "Побочный доход сохранен.")
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send additional income message:", slog.Any("error", err))
		}

	case "business_income":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Доход от бизнеса"
		if err := transaction.PostIncome(update, category, log); err != nil {
			log.Error("Failed to save business income:", slog.Any("error", err))
		}
		msg := tgbotapi.NewMessage(chatID, "Доход от бизнеса сохранен.")
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send business income message:", slog.Any("error", err))
		}

	case "investment_income":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Доход от инвестиций"
		if err := transaction.PostIncome(update, category, log); err != nil {
			log.Error("Failed to save investment income:", slog.Any("error", err))
		}
		msg := tgbotapi.NewMessage(chatID, "Доход от инвестиций сохранен.")
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send investment income message:", slog.Any("error", err))
		}

	case "state_payments":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Гос. выплаты"
		if err := transaction.PostIncome(update, category, log); err != nil {
			log.Error("Failed to save state payments:", slog.Any("error", err))
		}
		msg := tgbotapi.NewMessage(chatID, "Доход от государства сохранен.")
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send state income message:", slog.Any("error", err))
		}

	case "property_sales":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Продажа имущества"
		if err := transaction.PostIncome(update, category, log); err != nil {
			log.Error("Failed to save property sales:", slog.Any("error", err))
		}
		msg := tgbotapi.NewMessage(chatID, "Доход от продажи имущества сохранен.")
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send property sales message:", slog.Any("error", err))
		}

	case "other_income":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Прочие доходы"
		if err := transaction.PostIncome(update, category, log); err != nil {
			log.Error("Failed to save other income:", slog.Any("error", err))
		}
		msg := tgbotapi.NewMessage(chatID, "Прочие доходы сохранены.")
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send other income message:", slog.Any("error", err))
		}
	// expenses
	case "basic_expense":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Бытовые траты"
		if err := transaction.PostExpense(update, category, log); err != nil {
			log.Error("Failed to save basic expense:", slog.Any("error", err))
		}
		msg := tgbotapi.NewMessage(chatID, "Сумма базовых трат сохранена.")
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send basic expense message:", slog.Any("error", err))
		}

	case "regular_expense":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Регулярные платежи"
		if err := transaction.PostExpense(update, category, log); err != nil {
			log.Error("Failed to save regular expense:", slog.Any("error", err))
		}
		msg := tgbotapi.NewMessage(chatID, "Регулярный платеж сохранен.")
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send regular expense message:", slog.Any("error", err))
		}

	case "clothes":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Одежда"
		if err := transaction.PostExpense(update, category, log); err != nil {
			log.Error("Failed to save clothes:", slog.Any("error", err))
		}
		msg := tgbotapi.NewMessage(chatID, "Расход на обновление гардероба сохранен.")
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send clothes message:", slog.Any("error", err))
		}

	case "health":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Здоровье"
		if err := transaction.PostExpense(update, category, log); err != nil {
			log.Error("Failed to save health:", slog.Any("error", err))
		}
		msg := tgbotapi.NewMessage(chatID, "Расход на поддержание здоровья сохранен.")
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send health message:", slog.Any("error", err))
		}

	case "leisure_education":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Досуг и образование"
		if err := transaction.PostExpense(update, category, log); err != nil {
			log.Error("Failed to save leisure_education expense:", slog.Any("error", err))
		}
		msg := tgbotapi.NewMessage(chatID, "Расход сохранен.")
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send leisure_education message:", slog.Any("error", err))
		}

	case "investment_expense":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Инвестиции"
		if err := transaction.PostExpense(update, category, log); err != nil {
			log.Error("Failed to save investment expense:", slog.Any("error", err))
		}
		msg := tgbotapi.NewMessage(chatID, "Инвестиционный расход сохранен.")
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send investment expense message:", slog.Any("error", err))
		}

	case "other_expense":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Прочие расходы"
		if err := transaction.PostExpense(update, category, log); err != nil {
			log.Error("Failed to save other expense:", slog.Any("error", err))
		}
		msg := tgbotapi.NewMessage(chatID, "Прочие расходы сохранены.")
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send other expense message:", slog.Any("error", err))
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
			log.Error("Failed to send expense message:", slog.Any("error", err))
		}

	case "income":
		amount := update.Message.Text
		msg := tgbotapi.NewMessage(chatID, "Сумма прихода "+amount+" сохранена.")
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send income message:", slog.Any("error", err))
		}

	case "change_name":
		// Обновление имени пользователя в БД
		user := methodsForUser.UserMethod{}
		if err := user.UpdateUserName(update); err != nil {
			log.Error("Failed to update user name:", slog.Any("error", err))
		}
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s! Ваше имя успешно изменено.", update.Message.Text))
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send message about name change:", slog.Any("error", err))
		}

	case "change_currency":
		user := methodsForUser.UserMethod{}
		if err := user.UpdateUserCurrency(update); err != nil {
			log.Error("Failed to update user currency:", slog.Any("error", err))
		}
		msg := tgbotapi.NewMessage(chatID, "Ваша валюта изменена.")
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send message about currency change:", slog.Any("error", err))
		}
	}

	mu.Lock()
	delete(userStates, chatID) // удаляем состояние после обработки
	mu.Unlock()
}
