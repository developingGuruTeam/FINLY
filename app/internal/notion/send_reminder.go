package notion

import (
	"cachManagerApp/app/db/models"
	"cachManagerApp/database"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
	"time"
)

func StartReminderServiceWithCron(bot *tgbotapi.BotAPI) {
	c := cron.New()

	// TODO изменить на 12:00
	// "0 9 * * *" означает "каждый день в 9:00"
	c.AddFunc("*/1 * * * *", func() {
		err := processReminders(bot)
		if err != nil {
			fmt.Printf("Error while processing reminders: %v\n", err)
		}
	})

	// Запускаем cron-планировщик
	c.Start()
}

func processReminders(bot *tgbotapi.BotAPI) error {
	now := time.Now()

	var reminders []models.Reminder
	if err := database.DB.Where("next_reminder <= ?", now).Find(&reminders).Error; err != nil {
		return errors.New("Ошибка в поиске напоминаний")
	}

	for _, reminder := range reminders {
		sendReminder(bot, reminder)

		if reminder.Frequency == "неделя" {
			reminder.NextReminder = reminder.NextReminder.AddDate(0, 0, 7)
		} else if reminder.Frequency == "месяц" {
			reminder.NextReminder = reminder.NextReminder.AddDate(0, 1, 0)
		} else {
			// если нет частоты или одноразовое — можно либо удалить, либо больше не трогать ?
			reminder.NextReminder = time.Date(2100, 1, 1, 0, 0, 0, 0, time.Local)
		}

		if err := database.DB.Save(&reminder).Error; err != nil {
			log.Errorf("Error updating reminder: %v\n", err)
		}
	}
	return nil
}

func sendReminder(bot *tgbotapi.BotAPI, reminder models.Reminder) {
	chatID := int64(reminder.UserID)

	text := fmt.Sprintf(
		"☀️Мы к Вам с напоминанием\n\nНе забудьте оплатить сегодня платеж по категории '%s'\nСумма платежа: %d\n\n Хорошего дня 😊\n",
		reminder.Category,
		reminder.Amount,
	)

	msg := tgbotapi.NewMessage(chatID, text)
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Failed to send reminder to user %d: %v\n", chatID, err)
	}
}
