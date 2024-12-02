package TgBot

import (
	"cachManagerApp/app/internal/methodsForUser"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"sync"
)

type UserResponse struct {
	Action string `json:"action"`
}

var (
	userStates = make(map[int64]UserResponse) // мапа для хранения состояния пользователей
	mu         sync.Mutex                     // мьютекс для синхронизации доступа к мапе
)

// обработка нажатий на кнопки (команда приходит сюда)
func PushOnButton(bot *tgbotapi.BotAPI, update tgbotapi.Update, buttonCreator TelegramButtonCreator) {
	if update.Message != nil {
		// чат ID наполняется
		chatID := update.Message.Chat.ID
		mu.Lock()

		val, ok := userStates[chatID]
		mu.Unlock()
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
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "📥 Введите сумму прихода")
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
	}

	// Если команда или кнопка не обработаны, отправляем сообщение об ошибке
	if !handled {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🚫 Неизвестная команда. Повторите запрос.")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send unknown command message: %v", err)
		}
	}
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
	}

	mu.Lock()
	delete(userStates, chatID) // удаляем состояние после обработки
	mu.Unlock()
}
