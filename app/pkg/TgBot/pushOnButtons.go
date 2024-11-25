package TgBot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

// обработка нажатий на кнопки
func PushOnButton(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message != nil {
		switch update.Message.Text {
		case "Приход":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "📥 Введите сумму прихода")
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Failed to send message for income: %v", err)
			}

		case "Расход":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "📤 Введите сумму расхода")
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Failed to send message for expense: %v", err)
			}

		case "Отчеты":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "📊 Выберите тип отчета")
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Failed to send message for reports: %v", err)
			}

		case "Настройки":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "⚙ Выберите параметры")
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Failed to send message for settings: %v", err)
			}

		default:
			// защита от дурака
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🚫 Неизвестная команда. Повторите запрос")
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Failed to send message for unknown command: %v", err)
			}
		}
	}

	if update.CallbackQuery != nil {
		// обработка inline кнопок (пока две)
		switch update.CallbackQuery.Data {
		case "info":
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "Функциональность бота\n") // вот это по идее можно затереть, мне потестить не на чем
			if _, err := bot.Request(callback); err != nil {
				log.Printf("Failed to send callback response: %v", err)
			}
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Бот предназначен для:\n1. Ведения учета доходов и расходов\n2. Создания отчетов по различным критериям\n3. Экономического анализа")
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Failed to send message with bot functions: %v", err)
			}

		case "help":
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "Помощь по использованию бота") // вот это по идее можно затереть, мне потестить не на чем
			if _, err := bot.Request(callback); err != nil {
				log.Printf("Failed to send callback response: %v", err)
			}
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Команды бота:\n/info - Информация о боте") // дописать help, когда будет чем!
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Failed to send message with help info: %v", err)
			}
		}
	}
}
