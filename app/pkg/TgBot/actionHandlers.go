package TgBot

import (
	"cachManagerApp/app/internal/methodsForTransaction"
	"cachManagerApp/app/internal/methodsForUser"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Обработка транзакции
func handleTransactionAction(bot *tgbotapi.BotAPI, update tgbotapi.Update, transResp TransactionResponse) {
	chatID := update.Message.Chat.ID
	switch transResp.Action {
	// incomes
	case "salary":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Заработная плата"
		if err := transaction.PostIncome(update, category); err != nil {
			log.Printf("Failed to save salary: %v", err)
		}
		msg := tgbotapi.NewMessage(chatID, "Заработная плата сохранена.")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send salary message: %v", err)
		}

	case "additional_income":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Побочный доход"
		if err := transaction.PostIncome(update, category); err != nil {
			log.Printf("Failed to save additional income: %v", err)
		}
		msg := tgbotapi.NewMessage(chatID, "Побочный доход сохранен.")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send additional income message: %v", err)
		}

	case "business_income":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Доход от бизнеса"
		if err := transaction.PostIncome(update, category); err != nil {
			log.Printf("Failed to save business income: %v", err)
		}
		msg := tgbotapi.NewMessage(chatID, "Доход от бизнеса сохранен.")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send business income message: %v", err)
		}

	case "investment_income":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Доход от инвестиций"
		if err := transaction.PostIncome(update, category); err != nil {
			log.Printf("Failed to save investment income: %v", err)
		}
		msg := tgbotapi.NewMessage(chatID, "Доход от инвестиций сохранен.")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send investment income message: %v", err)
		}

	case "state_payments":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Гос.выплаты"
		if err := transaction.PostIncome(update, category); err != nil {
			log.Printf("Failed to save investment income: %v", err)
		}
		msg := tgbotapi.NewMessage(chatID, "Доход от государства сохранен.")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send state income message: %v", err)
		}

	case "property_sales":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Доход от продажи имущества"
		if err := transaction.PostIncome(update, category); err != nil {
			log.Printf("Failed to save investment income: %v", err)
		}
		msg := tgbotapi.NewMessage(chatID, "Доход от продажи имущества сохранен.")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send investment income message: %v", err)
		}

	case "other_income":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Прочие доходы"
		if err := transaction.PostIncome(update, category); err != nil {
			log.Printf("Failed to save other income: %v", err)
		}
		msg := tgbotapi.NewMessage(chatID, "Прочие доходы сохранены.")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send other income message: %v", err)
		}
	// expenses
	case "basic_expense":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Бытовые траты"
		if err := transaction.PostExpense(update, category); err != nil {
			log.Printf("Failed to save basic expense: %v", err)
		}
		msg := tgbotapi.NewMessage(chatID, "Сумма базовых трат сохранена.")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send basic expense message: %v", err)
		}

	case "regular_expense":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Побочный доход"
		if err := transaction.PostExpense(update, category); err != nil {
			log.Printf("Failed to save regular expense: %v", err)
		}
		msg := tgbotapi.NewMessage(chatID, "Регулярный платеж сохранен.")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send regular expense message: %v", err)
		}

	case "clothes":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Одежда"
		if err := transaction.PostExpense(update, category); err != nil {
			log.Printf("Failed to save clothes: %v", err)
		}
		msg := tgbotapi.NewMessage(chatID, "Расход на обновление гардероба сохранен.")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send clothes message: %v", err)
		}

	case "health":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Здоровье"
		if err := transaction.PostExpense(update, category); err != nil {
			log.Printf("Failed to save health: %v", err)
		}
		msg := tgbotapi.NewMessage(chatID, "Расход на поддержание здоровья сохранен.")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send health message: %v", err)
		}

	case "leisure_education":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Досуг и образование"
		if err := transaction.PostExpense(update, category); err != nil {
			log.Printf("Failed to save leisure_education expense: %v", err)
		}
		msg := tgbotapi.NewMessage(chatID, "Расход сохранен.")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send leisure_education message: %v", err)
		}

	case "investment_expense":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Расход на инвестиции"
		if err := transaction.PostExpense(update, category); err != nil {
			log.Printf("Failed to save investment expense: %v", err)
		}
		msg := tgbotapi.NewMessage(chatID, "Инвестиционный расход сохранен.")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send investment expense message: %v", err)
		}

	case "other_expense":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Прочие доходы"
		if err := transaction.PostExpense(update, category); err != nil {
			log.Printf("Failed to save other expense: %v", err)
		}
		msg := tgbotapi.NewMessage(chatID, "Прочие расходы сохранены.")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send other expense message: %v", err)
		}
	}
	mu.Lock()
	delete(transactionStates, chatID) // удаляем состояние после обработки
	mu.Unlock()

}

func handleUserAction(bot *tgbotapi.BotAPI, update tgbotapi.Update, userResp UserResponse) {
	chatID := update.Message.Chat.ID

	switch userResp.Action {
	case "expense":
		amount := update.Message.Text
		msg := tgbotapi.NewMessage(chatID, "Сумма расхода "+amount+" сохранена.")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Ошибка отправки сообщения о расходе: %v", err)
		}

	case "income":
		amount := update.Message.Text
		msg := tgbotapi.NewMessage(chatID, "Сумма прихода "+amount+" сохранена.")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Ошибка отправки сообщения о приходе: %v", err)
		}

	case "change_name":
		// Обновление имени пользователя в БД
		user := methodsForUser.UserMethod{}
		if err := user.UpdateUserName(update); err != nil {
			log.Errorf("Ошибка обновления имени пользователя: %v", err)
		}
		msg := tgbotapi.NewMessage(chatID, "Ваше имя успешно изменено.")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Ошибка отправки сообщения об изменении имени: %v", err)
		}

	case "change_currency":
		user := methodsForUser.UserMethod{}
		if err := user.UpdateUserCurrency(update); err != nil {
			log.Errorf("Ошибка обновления валюты: %v", err)
		}
		msg := tgbotapi.NewMessage(chatID, "Ваша валюта изменена.")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Ошибка отправки сообщения об изменении валюты: %v", err)
		}
	}

	mu.Lock()
	delete(userStates, chatID) // удаляем состояние после обработки
	mu.Unlock()
}
