package notion

import (
	"cachManagerApp/app/db/models"
	"cachManagerApp/app/internal/notion/rulesForNotion"
	"cachManagerApp/app/pkg/ButtonsCreate"
	"cachManagerApp/app/pkg/logger"
	"cachManagerApp/database"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
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

		msg := tgbotapi.NewMessage(chatID, "Введите дату следующего платежа (формат: ДД.ММ.ГГГГ)")
		_, err := bot.Send(msg)
		if err != nil {
			log.Errorf("Ошибка в отправке сообщения в категории напоминания %v", err)
		}
		return

	case reminder.NextReminder.IsZero():
		// Проверяем и сохраняем дату платежа
		nextReminder, err := rulesForNotion.ValidateRightTime(update.Message.Text)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, err.Error())
			_, _ = bot.Send(msg)
			return
		}

		reminder.NextReminder = nextReminder

		msg := tgbotapi.NewMessage(chatID, "Введите сумму платежа (только цифры):")
		_, err = bot.Send(msg)
		if err != nil {
			log.Errorf("Ошибка в отправке сообщения суммы: %v", err)
		}

	case reminder.Amount == 0:
		// Получаем сумму платежа
		amount, err := strconv.Atoi(update.Message.Text)
		if err != nil || amount <= 0 {
			msg := tgbotapi.NewMessage(chatID, "Неверный ввод. Пожалуйста, введите положительное целое число.")
			_, _ = bot.Send(msg)
			return
		}

		reminder.Amount = amount
		reminder.CreatedAt = time.Now()

		if err := database.DB.Create(&reminder).Error; err != nil {
			log.Errorf("Ошибка при сохранении напоминания: %v", err)
			msg := tgbotapi.NewMessage(chatID, "Произошла ошибка при сохранении напоминания. Попробуйте позже.")
			_, _ = bot.Send(msg)
			return
		}

		menuMain := ButtonsCreate.TelegramButtonCreator{}
		back := menuMain.CreateMainMenuButtons()
		msg := tgbotapi.NewMessage(chatID, "Напоминание успешно создано!")

		fmt.Println(reminder)

		msg.ReplyMarkup = back
		if _, err := bot.Send(msg); err != nil {
			log.Errorf("Failed to send main menu: %v", err)
		}

		// Удаляем напоминание из состояния
		delete(RemindersStates, chatID)
		return
	}
}
