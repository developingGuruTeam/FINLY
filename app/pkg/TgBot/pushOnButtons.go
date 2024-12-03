package TgBot

import (
	"cachManagerApp/app/internal/methodsForTransaction"
	"cachManagerApp/app/internal/methodsForUser"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"sync"
)

type TransactionResponse struct {
	Action string `json:"action"`
}

type UserResponse struct {
	Action string `json:"action"`
}

var (
	userStates        = make(map[int64]UserResponse)        // мапа для хранения состояния пользователей
	mu                sync.Mutex                            // мьютекс для синхронизации доступа к мапе
	transactionStates = make(map[int64]TransactionResponse) // мапа для хранения состояния транзакций
)

// обработка нажатий на кнопки (команда приходит сюда)
func PushOnButton(bot *tgbotapi.BotAPI, update tgbotapi.Update, buttonCreator TelegramButtonCreator) {
	if update.Message != nil {
		// чат ID наполняется
		chatID := update.Message.Chat.ID
		mu.Lock()
		val2, ok2 := transactionStates[chatID]
		val, ok := userStates[chatID]
		mu.Unlock()

		if ok2 && val2.Action != "" {
			handleTransactionAction(bot, update, val2)
			return
		}

		// если в ней лежит ключ, то переходит к действию, если нет, то ждет отклика
		if ok && val.Action != "" {
			handleUserAction(bot, update, val)
			return
		}
		handleButtonPress(bot, update, buttonCreator)
	}
}

func handleButtonPress(bot *tgbotapi.BotAPI, update tgbotapi.Update, buttonCreator TelegramButtonCreator) {
	chatID := update.Message.Chat.ID
	handled := false
	switch update.Message.Text {

	// ОПИСАНИЕ КНОПОК МЕНЮ
	case "📥 Приход":
		incomeMenu := buttonCreator.CreateIncomeMenuButtons()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "⚙ Выберите категорию")
		msg.ReplyMarkup = incomeMenu
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send message for income: %v", err)
		}
		handled = true

	case "📤 Расход":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "📤 Введите сумму расхода")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send message for expense: %v", err)
		}
		handled = true

	case "📊 Отчеты":
		reportMenu := buttonCreator.CreateReportsMenuButtons()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "📊 Выберите тип отчета")
		msg.ReplyMarkup = reportMenu
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send message for reports: %v", err)
		}
		handled = true

	case "⚙ Настройки":
		settingsMenu := buttonCreator.CreateSettingsMenuButtons()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "⚙ Выберите параметры")
		msg.ReplyMarkup = settingsMenu
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send message for settings: %v", err)
		}
		handled = true

	case "⬅ В меню":
		mainMenu := buttonCreator.CreateMainMenuButtons()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вы вернулись в главное меню")
		msg.ReplyMarkup = mainMenu
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send main menu: %v", err)
		}
		handled = true

	// ОПИСАНИЕ ИНЛАЙН КОММАНД
	case "/info":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "📍 Бот предназначен для:\n ▪ Ведения учета доходов и расходов\n ▪ Создания отчетов по различным критериям\n ▪ Экономического анализа")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send /info message: %v", err)
		}
		handled = true

	case "/help":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "📌 Команды бота:\n/info - Информация о боте\n/help - Помощь по использованию бота") // дописать нормальный хэлп
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send /help message: %v", err)
		}
		handled = true

	case "/hi":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, randomTextForHi()) // дописать нормальный хэлп
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send /help message: %v", err)
		}
		handled = true

	// подумать над состоянием ответа ТГ
	case "🎭 Изменить имя":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите Ваше новое имя")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send /help message: %v", err)
		}
		mu.Lock()
		userStates[chatID] = UserResponse{Action: "change_name"}
		mu.Unlock()
		handled = true

	case "💱 Изменить валюту":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите валюту")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send /help message: %v", err)
		}
		mu.Lock()
		userStates[chatID] = UserResponse{Action: "change_currency"}
		mu.Unlock()
		handled = true

		// приходы

	case "📥 Заработная плата":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите сумму заработной платы.\nЧерез запятую можно добавить комментарий")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send /help message: %v", err)
		}
		mu.Lock()
		transactionStates[chatID] = TransactionResponse{Action: "salary"}
		mu.Unlock()
		handled = true

	case "📤 Дополнительный доход":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите сумму дополнительного дохода.\nЧерез запятую можно добавить комментарий")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send /help message: %v", err)
		}
		mu.Lock()
		transactionStates[chatID] = TransactionResponse{Action: "additional_income"}
		mu.Unlock()
		handled = true

	case "📥 Доход от бизнеса":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите сумму дохода от бизнеса.\nЧерез запятую можно добавить комментарий")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send /help message: %v", err)
		}
		mu.Lock()
		transactionStates[chatID] = TransactionResponse{Action: "business_income"}
		mu.Unlock()
		handled = true

	case "📥 Доход от инвестиций":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите сумму дохода от инвестиций.\nЧерез запятую можно добавить комментарий")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send /help message: %v", err)
		}
		mu.Lock()
		transactionStates[chatID] = TransactionResponse{Action: "investment_income"}
		mu.Unlock()
		handled = true

	case "📥 Государственные выплаты":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите сумму государственных выплат.\nЧерез запятую можно добавить комментарий")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send /help message: %v", err)
		}
		mu.Lock()
		transactionStates[chatID] = TransactionResponse{Action: "state_payments"}
		mu.Unlock()
		handled = true

	case "📤 Продажа имущества":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите сумму продажи имущества.\nЧерез запятую можно добавить комментарий")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send /help message: %v", err)
		}
		mu.Lock()
		transactionStates[chatID] = TransactionResponse{Action: "property_sales"}
		mu.Unlock()
		handled = true

	case "📥 Прочее":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите сумму прочего дохода.\nЧерез запятую можно добавить комментарий")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send /help message: %v", err)
		}
		mu.Lock()
		transactionStates[chatID] = TransactionResponse{Action: "other_income"}
		mu.Unlock()
		handled = true
	}

	// Если команда или кнопка не обработаны, отправляем сообщение об ошибке
	if !handled {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🚫 Неизвестная команда. Повторите запрос.")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send unknown command message: %v", err)
		}
	}
}

func handleTransactionAction(bot *tgbotapi.BotAPI, update tgbotapi.Update, transResp TransactionResponse) {
	chatID := update.Message.Chat.ID

	switch transResp.Action {
	case "salary":
		transaction := methodsForTransaction.TransactionsMethod{}
		category := "Зарплата"
		if err := transaction.PostIncome(update, category); err != nil {
			log.Printf("Failed to save salary: %v", err)
		}
		msg := tgbotapi.NewMessage(chatID, "Зарплата сохранена.")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send salary message: %v", err)
		}

	case "additional_income":
		trasaction := methodsForTransaction.TransactionsMethod{}
		category := "Дополнительный доход"
		if err := trasaction.PostIncome(update, category); err != nil {
			log.Printf("Failed to save additional income: %v", err)
		}
		msg := tgbotapi.NewMessage(chatID, "Дополнительный доход сохранен.")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send additional income message: %v", err)
		}

	case "business_income":
		trasaction := methodsForTransaction.TransactionsMethod{}
		category := "Доход от бизнеса"
		if err := trasaction.PostIncome(update, category); err != nil {
			log.Printf("Failed to save business income: %v", err)
		}
		msg := tgbotapi.NewMessage(chatID, "Доход от бизнеса сохранен.")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send business income message: %v", err)
		}

	case "investment_income":
		trasaction := methodsForTransaction.TransactionsMethod{}
		category := "Доход от инвестиций"
		if err := trasaction.PostIncome(update, category); err != nil {
			log.Printf("Failed to save investment income: %v", err)
		}
		msg := tgbotapi.NewMessage(chatID, "Доход от инвестиций сохранен.")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send investment income message: %v", err)
		}

	case "other_income":
		trasaction := methodsForTransaction.TransactionsMethod{}
		category := "Прочие доходы"
		if err := trasaction.PostIncome(update, category); err != nil {
			log.Printf("Failed to save other income: %v", err)
		}
		msg := tgbotapi.NewMessage(chatID, "Прочие доходы сохранены.")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send other income message: %v", err)
		}
	}
	mu.Lock()
	delete(transactionStates, chatID) // удаляем состояние после обработки
	mu.Unlock()

}

func handleUserAction(bot *tgbotapi.BotAPI, update tgbotapi.Update, userResp UserResponse) {
	chatID := update.Message.Chat.ID

	switch userResp.Action {
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
			log.Printf("Ошибка обновления имени пользователя: %v", err)
			return
		}
		msg := tgbotapi.NewMessage(chatID, "Ваше имя успешно изменено.")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Ошибка отправки сообщения об изменении имени: %v", err)
		}

	case "change_currency":
		user := methodsForUser.UserMethod{}
		if err := user.UpdateUserCurrency(update); err != nil {
			log.Printf("Ошибка обновления валюты: %v", err)
			return
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
