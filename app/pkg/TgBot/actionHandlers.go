package TgBot

import (
	"cachManagerApp/app/internal/methodsForTransaction"
	"cachManagerApp/app/pkg/ButtonsCreate"
	"log/slog"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Обработка транзакции
func handleTransactionAction(bot *tgbotapi.BotAPI, update tgbotapi.Update, transResp TransactionResponse, buttonCreator ButtonsCreate.TelegramButtonCreator, log *slog.Logger) {
	chatID := update.Message.Chat.ID
	switch transResp.Action {
	// Доходы
	case "salary":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Заработная плата"
		sum, err := strconv.Atoi(update.Message.Text)
		if err != nil || sum <= 0 {
			msg := tgbotapi.NewMessage(chatID, "🚫 Введите корректное положительное целое число.")
			bot.Send(msg)
			return
		}

		if err := transaction.PostIncome(update, category, log); err != nil {
			log.Info("Failed to save salary: %s", log.With("error", err))
			msg := tgbotapi.NewMessage(chatID, "❌ Ошибка сохранения транзакции.")
			bot.Send(msg)
			return
		}

		doneMsg := "✅ Заработная плата успешно сохранена."
		returnToMainMenu(bot, chatID, buttonCreator, doneMsg)

	case "additional_income":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Побочный доход"
		sum, err := strconv.Atoi(update.Message.Text)
		if err != nil || sum <= 0 {
			msg := tgbotapi.NewMessage(chatID, "🚫 Введите корректное положительное целое число.")
			bot.Send(msg)
			return
		}

		if err := transaction.PostIncome(update, category, log); err != nil {
			log.Info("Failed to save additional income: %s", log.With("error", err))
			msg := tgbotapi.NewMessage(chatID, "❌ Ошибка сохранения транзакции.")
			bot.Send(msg)
			return
		}

		doneMsg := "✅ Побочный доход успешно сохранен."
		returnToMainMenu(bot, chatID, buttonCreator, doneMsg)

	case "business_income":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Доход от бизнеса"
		sum, err := strconv.Atoi(update.Message.Text)
		if err != nil || sum <= 0 {
			msg := tgbotapi.NewMessage(chatID, "🚫 Введите корректное положительное целое число.")
			bot.Send(msg)
			return
		}

		if err := transaction.PostIncome(update, category, log); err != nil {
			log.Info("Failed to save business income: %s", log.With("error", err))
			msg := tgbotapi.NewMessage(chatID, "❌ Ошибка сохранения транзакции.")
			bot.Send(msg)
			return
		}

		doneMsg := "✅ Доход от бизнеса успешно сохранен."
		returnToMainMenu(bot, chatID, buttonCreator, doneMsg)

	case "investment_income":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Доход от инвестиций"
		sum, err := strconv.Atoi(update.Message.Text)
		if err != nil || sum <= 0 {
			msg := tgbotapi.NewMessage(chatID, "🚫 Введите корректное положительное целое число.")
			bot.Send(msg)
			return
		}

		if err := transaction.PostIncome(update, category, log); err != nil {
			log.Info("Failed to save investment income: %s", log.With("error", err))
			msg := tgbotapi.NewMessage(chatID, "❌ Ошибка сохранения транзакции.")
			bot.Send(msg)
			return
		}

		doneMsg := "✅ Доход от инвестиций успешно сохранен."
		returnToMainMenu(bot, chatID, buttonCreator, doneMsg)

	case "state_payments":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Гос. выплаты"
		sum, err := strconv.Atoi(update.Message.Text)
		if err != nil || sum <= 0 {
			msg := tgbotapi.NewMessage(chatID, "🚫 Введите корректное положительное целое число.")
			bot.Send(msg)
			return
		}

		if err := transaction.PostIncome(update, category, log); err != nil {
			log.Info("Failed to save state payments: %s", log.With("error", err))
			msg := tgbotapi.NewMessage(chatID, "❌ Ошибка сохранения транзакции.")
			bot.Send(msg)
			return
		}

		doneMsg := "✅ Гос. выплаты успешно сохранены."
		returnToMainMenu(bot, chatID, buttonCreator, doneMsg)

	case "property_sales":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Продажа имущества"
		sum, err := strconv.Atoi(update.Message.Text)
		if err != nil || sum <= 0 {
			msg := tgbotapi.NewMessage(chatID, "🚫 Введите корректное положительное целое число.")
			bot.Send(msg)
			return
		}

		if err := transaction.PostIncome(update, category, log); err != nil {
			log.Info("Failed to save property sales: %s", log.With("error", err))
			msg := tgbotapi.NewMessage(chatID, "❌ Ошибка сохранения транзакции.")
			bot.Send(msg)
			return
		}

		doneMsg := "✅ Продажа имущества успешно сохранена."
		returnToMainMenu(bot, chatID, buttonCreator, doneMsg)

	case "other_income":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Прочие доходы"
		sum, err := strconv.Atoi(update.Message.Text)
		if err != nil || sum <= 0 {
			msg := tgbotapi.NewMessage(chatID, "🚫 Введите корректное положительное целое число.")
			bot.Send(msg)
			return
		}

		if err := transaction.PostIncome(update, category, log); err != nil {
			log.Info("Failed to save other income: %s", log.With("error", err))
			msg := tgbotapi.NewMessage(chatID, "❌ Ошибка сохранения транзакции.")
			bot.Send(msg)
			return
		}

		doneMsg := "✅ Прочие доходы успешно сохранены."
		returnToMainMenu(bot, chatID, buttonCreator, doneMsg)

	// Расходы
	case "basic_expense":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Бытовые траты"
		sum, err := strconv.Atoi(update.Message.Text)
		if err != nil || sum <= 0 {
			msg := tgbotapi.NewMessage(chatID, "🚫 Введите корректное положительное целое число.")
			bot.Send(msg)
			return
		}

		if err := transaction.PostExpense(update, category, log); err != nil {
			log.Info("Failed to save basic expense: %s", log.With("error", err))
			msg := tgbotapi.NewMessage(chatID, "❌ Ошибка сохранения транзакции.")
			bot.Send(msg)
			return
		}

		doneMsg := "✅ Бытовые траты успешно сохранены."
		returnToMainMenu(bot, chatID, buttonCreator, doneMsg)

	case "regular_expense":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Регулярные платежи"
		sum, err := strconv.Atoi(update.Message.Text)
		if err != nil || sum <= 0 {
			msg := tgbotapi.NewMessage(chatID, "🚫 Введите корректное положительное целое число.")
			bot.Send(msg)
			return
		}

		if err := transaction.PostExpense(update, category, log); err != nil {
			log.Info("Failed to save regular expense: %s", log.With("error", err))
			msg := tgbotapi.NewMessage(chatID, "❌ Ошибка сохранения транзакции.")
			bot.Send(msg)
			return
		}

		doneMsg := "✅ Регулярные платежи успешно сохранены."
		returnToMainMenu(bot, chatID, buttonCreator, doneMsg)

	case "clothes":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Одежда"
		sum, err := strconv.Atoi(update.Message.Text)
		if err != nil || sum <= 0 {
			msg := tgbotapi.NewMessage(chatID, "🚫 Введите корректное положительное целое число.")
			bot.Send(msg)
			return
		}

		if err := transaction.PostExpense(update, category, log); err != nil {
			log.Info("Failed to save clothes: %s", log.With("error", err))
			msg := tgbotapi.NewMessage(chatID, "❌ Ошибка сохранения транзакции.")
			bot.Send(msg)
			return
		}

		doneMsg := "✅ Траты на одежду успешно сохранены."
		returnToMainMenu(bot, chatID, buttonCreator, doneMsg)

	case "health":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Здоровье"
		sum, err := strconv.Atoi(update.Message.Text)
		if err != nil || sum <= 0 {
			msg := tgbotapi.NewMessage(chatID, "🚫 Введите корректное положительное целое число.")
			bot.Send(msg)
			return
		}

		if err := transaction.PostExpense(update, category, log); err != nil {
			log.Info("Failed to save health: %s", log.With("error", err))
			msg := tgbotapi.NewMessage(chatID, "❌ Ошибка сохранения транзакции.")
			bot.Send(msg)
			return
		}

		doneMsg := "✅ Траты на здоровье успешно сохранены."
		returnToMainMenu(bot, chatID, buttonCreator, doneMsg)

	case "leisure_education":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Досуг и образование"
		sum, err := strconv.Atoi(update.Message.Text)
		if err != nil || sum <= 0 {
			msg := tgbotapi.NewMessage(chatID, "🚫 Введите корректное положительное целое число.")
			bot.Send(msg)
			return
		}

		if err := transaction.PostExpense(update, category, log); err != nil {
			log.Info("Failed to save leisure and education expense: %s", log.With("error", err))
			msg := tgbotapi.NewMessage(chatID, "❌ Ошибка сохранения транзакции.")
			bot.Send(msg)
			return
		}

		doneMsg := "✅ Траты на досуг и образование успешно сохранены."
		returnToMainMenu(bot, chatID, buttonCreator, doneMsg)

	case "investment_expense":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Инвестиции"
		sum, err := strconv.Atoi(update.Message.Text)
		if err != nil || sum <= 0 {
			msg := tgbotapi.NewMessage(chatID, "🚫 Введите корректное положительное целое число.")
			bot.Send(msg)
			return
		}

		if err := transaction.PostExpense(update, category, log); err != nil {
			log.Info("Failed to save investment expense: %s", log.With("error", err))
			msg := tgbotapi.NewMessage(chatID, "❌ Ошибка сохранения транзакции.")
			bot.Send(msg)
			return
		}

		doneMsg := "✅ Инвестиционные траты успешно сохранены."
		returnToMainMenu(bot, chatID, buttonCreator, doneMsg)

	case "other_expense":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Прочие расходы"
		sum, err := strconv.Atoi(update.Message.Text)
		if err != nil || sum <= 0 {
			msg := tgbotapi.NewMessage(chatID, "🚫 Введите корректное положительное целое число.")
			bot.Send(msg)
			return
		}

		if err := transaction.PostExpense(update, category, log); err != nil {
			log.Info("Failed to save other expenses: %s", log.With("error", err))
			msg := tgbotapi.NewMessage(chatID, "❌ Ошибка сохранения транзакции.")
			bot.Send(msg)
			return
		}

		doneMsg := "✅ Прочие траты успешно сохранены."
		returnToMainMenu(bot, chatID, buttonCreator, doneMsg)
	}
	mu.Lock()
	delete(transactionStates, chatID) // удаляем состояние после обработки
	mu.Unlock()

}
