package TgBot

import (
	"cachManagerApp/app/internal/methodsForUser"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

// обработка нажатий на кнопки
func PushOnButton(bot *tgbotapi.BotAPI, update tgbotapi.Update, buttonCreator TelegramButtonCreator) {
	if update.Message != nil {
		handled := false // флаг, чтобы понимать обработана ли команда/кнопка

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

		case "🎭 Изменить имя":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите Ваше новое имя")
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Failed to send /help message: %v", err)
			}
			user := methodsForUser.UserMethod{}
			user.WaitingUpdate = true

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
}
