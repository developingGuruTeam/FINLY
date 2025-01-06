package TgBot

import (
	"cachManagerApp/app/internal/methodsForTransaction"
	"cachManagerApp/app/internal/methodsForUser"
	"cachManagerApp/app/pkg/ButtonsCreate"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

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

func handleUserAction(bot *tgbotapi.BotAPI, update tgbotapi.Update, userResp UserResponse, buttonCreator ButtonsCreate.TelegramButtonCreator, log *slog.Logger) {
	chatID := update.Message.Chat.ID

	switch userResp.Action {

	case "change_name":
		newName := strings.TrimSpace(update.Message.Text) // убираем пробелы по обе стороны, если есть

		// проверка нового имени: только буквы и длина от 1 до 32 символов
		var validName bool = true
		for _, symbol := range newName {
			if !unicode.IsLetter(symbol) && symbol != ' ' { // имя только из букв и пробелов
				validName = false
				break
			}
		}

		if utf8.RuneCountInString(newName) == 0 || utf8.RuneCountInString(newName) > 32 || !validName {
			msg := tgbotapi.NewMessage(chatID, "🚫 Некорректное имя. Имя должно содержать только буквы и быть не более 32 символов.")
			bot.Send(msg)
			return
		}

		// обновляем имя пользователя
		user := methodsForUser.UserMethod{}
		if err := user.UpdateUserName(update); err != nil {
			log.Error("Ошибка обновления имени пользователя", log.With("error", err))
			msg := tgbotapi.NewMessage(chatID, "❌ Ошибка при обновлении имени.")
			bot.Send(msg)
			return
		}

		msgDone := fmt.Sprintf("✅ Ваше имя успешно изменено на %s.", newName)
		returnToMainMenu(bot, chatID, buttonCreator, msgDone)

	case "change_currency":
		newCurrency := strings.ToLower(update.Message.Text) // преобразуем в нижний регистр
		// проверка новой валюты на алфавит
		var okCurrency bool = true
		for _, symbol := range newCurrency {
			if !unicode.IsLetter(symbol) {
				okCurrency = false
				break
			}
		}
		// проверка валюты на длину
		if utf8.RuneCountInString(newCurrency) != 3 || okCurrency != true {
			msg := tgbotapi.NewMessage(chatID, "🚫 Некорректный формат валюты. Валюта должна содержать только буквы и быть не более 3 символов.")
			bot.Send(msg)
			return
		}
		// обновляем валюту пользователя
		user := methodsForUser.UserMethod{}
		if err := user.UpdateUserCurrency(update); err != nil {
			log.Error("Ошибка обновления валюты", log.With("error", err))
			msg := tgbotapi.NewMessage(chatID, "❌ Ошибка при обновлении валюты")
			bot.Send(msg)
			return
		}

		msgDone := fmt.Sprintf("✅ Ваша валюта изменена на %s.", newCurrency)
		returnToMainMenu(bot, chatID, buttonCreator, msgDone)
	}

	mu.Lock()
	delete(userStates, chatID) // удаляем состояние после обработки
	mu.Unlock()
}

// возврат кнопок меню и удаления состояния после обработки транзакции
func returnToMainMenu(bot *tgbotapi.BotAPI, chatID int64, buttonCreator ButtonsCreate.TelegramButtonCreator, msg string) {
	// создаем кнопки главного меню
	mainMenu := buttonCreator.CreateMainMenuButtons()

	// отправляем пустое сообщение с кнопками
	menuMsg := tgbotapi.NewMessage(chatID, msg)
	menuMsg.ReplyMarkup = mainMenu // показываем кнопки
	bot.Send(menuMsg)

	delete(transactionStates, chatID)
}
