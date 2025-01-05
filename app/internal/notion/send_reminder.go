package notion

import (
	"cachManagerApp/app/db/models"
	"cachManagerApp/database"
	"errors"
	"fmt"
	"log/slog"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
)

func StartReminderServiceWithCron(bot *tgbotapi.BotAPI, log *slog.Logger) {
	c := cron.New()

	// TODO изменить на 12:00
	// "0 9 * * *" означает "каждый день в 9:00"
	c.AddFunc("*/1 * * * *", func() {
		err := processReminders(bot, log)
		if err != nil {
			log.Error(
				"Error while processing reminders: %v\n",
				log.With("error", err),
			)
		}
	})

	// Запускаем cron-планировщик
	c.Start()
}

func processReminders(bot *tgbotapi.BotAPI, log *slog.Logger) error {
	now := time.Now()

	var reminders []models.Reminder
	if err := database.DB.Where("next_reminder <= ?", now).Find(&reminders).Error; err != nil {
		return errors.New("Ошибка в поиске напоминаний")
	}

	for _, reminder := range reminders {
		sendReminder(bot, reminder, log)

		if reminder.Frequency == "неделя" {
			reminder.NextReminder = reminder.NextReminder.AddDate(0, 0, 7)
		} else if reminder.Frequency == "месяц" {
			reminder.NextReminder = reminder.NextReminder.AddDate(0, 1, 0)
		} else {
			// если нет частоты или одноразовое — можно либо удалить, либо больше не трогать ?
			reminder.NextReminder = time.Date(2100, 1, 1, 0, 0, 0, 0, time.Local)
		}

		if err := database.DB.Save(&reminder).Error; err != nil {
			log.Error("Error updating reminder:", log.With("error", err))
		}
	}
	return nil
}

func sendReminder(bot *tgbotapi.BotAPI, reminder models.Reminder, log *slog.Logger) {
	chatID := int64(reminder.UserID)

	text := fmt.Sprintf(
		"☀️Мы к Вам с напоминанием\n\nНе забудьте оплатить сегодня платеж по категории '%s'\nСумма платежа: %d\n\n Хорошего дня 😊\n",
		reminder.Category,
		reminder.Amount,
	)

	msg := tgbotapi.NewMessage(chatID, text)
	_, err := bot.Send(msg)
	if err != nil {
		log.Error(
			"Failed to send reminder to user %d:",
			log.With("user_id", chatID),
			log.With("error", err),
		)
	}
}
