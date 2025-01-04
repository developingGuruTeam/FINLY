package notion

import (
	"cachManagerApp/app/db/models"
	"cachManagerApp/app/pkg/ButtonsCreate"
	"cachManagerApp/app/pkg/logger"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

var log = logger.GetLogger()
var RemindersStates = map[int64]*models.Reminder{}

func StartReminder(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID

	// создаем новый Reminder для пользователя
	RemindersStates[chatID] = &models.Reminder{
		UserID: uint64(chatID),
	}
}

func HandleReminderInput(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	reminder, exists := RemindersStates[chatID]
	if !exists {
		return
	}

	switch {
	case reminder.Frequency == "":
		// получаем частоту платежа

		if update.Message.Text == "⬅ В меню" {
			menuMain := ButtonsCreate.TelegramButtonCreator{}
			back := menuMain.CreateMainMenuButtons()
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вы вернулись в главное меню")
			msg.ReplyMarkup = back
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Failed to send main menu: %v", err)
			}
			// Удаляем напоминание из состояния
			delete(RemindersStates, update.Message.Chat.ID)
			return
		}

		frequency := update.Message.Text
		if frequency != "🫠 Через неделю" && frequency != "🌙 Через месяц" {
			msg := tgbotapi.NewMessage(chatID, "Неверный ввод. Пожалуйста, выберите '🫠 Через неделю' или '🌙 Через месяц'.")
			_, _ = bot.Send(msg)
			return
		}

		if frequency == "🫠 Через неделю" {
			reminder.Frequency = "неделя"
		}
		if frequency == "🌙 Через месяц" {
			reminder.Frequency = "месяц"
		}

		// Переходим к следующему этапу — названию платежа
		msg := tgbotapi.NewMessage(chatID, "Введите название следующего платежа:")
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true) // Убираем кнопки
		_, _ = bot.Send(msg)

	case reminder.Category == "":
		// Получаем название платежа
		reminder.Category = update.Message.Text

		msg := tgbotapi.NewMessage(chatID, "Введите дату следующего платежа (формат: ГГГГ-ММ-ДД):")
		_, err := bot.Send(msg)
		if err != nil {
			log.Errorf("Ошибка в отправке сообщения в категории напоминания %v", err)
		}
		return

	case reminder.NextReminder.IsZero():
		// Проверяем и сохраняем дату платежа
		date, err := time.Parse("2006-01-02", update.Message.Text)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Неверный формат даты. Попробуйте ещё раз (ГГГГ-ММ-ДД):")
			_, _ = bot.Send(msg)
			return
		}

		reminder.NextReminder = date

		menuMain := ButtonsCreate.TelegramButtonCreator{}
		back := menuMain.CreateMainMenuButtons()
		msg := tgbotapi.NewMessage(chatID, "Напоминание успешно создано!")

		fmt.Println(reminder)

		msg.ReplyMarkup = back
		if _, err := bot.Send(msg); err != nil {
			log.Errorf("Failed to send main menu: %v", err)
		}
		// Удаляем напоминание из состояния
		delete(RemindersStates, update.Message.Chat.ID)
		return
	}
}
