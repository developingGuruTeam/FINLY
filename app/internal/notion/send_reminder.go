package notion

import (
	"cachManagerApp/app/db/models"
	"cachManagerApp/database"
	"fmt"
	"log/slog"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
)

func StartReminderServiceWithCron(bot *tgbotapi.BotAPI, log *slog.Logger) {
	c := cron.New()

	// "0 11 * * *" означает "каждый день в 11:00"
	c.AddFunc("0 10 * * *", func() { // использует локально время сервера
		err := processReminders(bot, log)
		if err != nil {
			log.Error(
				"Error while processing reminders: %v\n",
				"error", err)
		}
	})

	// Запускаем cron-планировщик
	c.Start()
}

func processReminders(bot *tgbotapi.BotAPI, log *slog.Logger) error {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	var reminders []models.Reminder
	if err := database.DB.Where("next_reminder BETWEEN ? AND ?", startOfDay, endOfDay).Find(&reminders).Error; err != nil {
		return fmt.Errorf("Ошибка в поиске напоминаний :%v", err)
	}

	remindersByUser := make(map[uint64][]models.Reminder)
	for _, reminder := range reminders {
		remindersByUser[reminder.UserID] = append(remindersByUser[reminder.UserID], reminder)
	}

	for userID, userReminders := range remindersByUser {
		sendReminders(bot, userID, userReminders, log)

		for _, reminder := range userReminders {
			if reminder.Frequency == "неделя" {
				reminder.NextReminder = reminder.NextReminder.AddDate(0, 0, 7)
			} else if reminder.Frequency == "месяц" {
				reminder.NextReminder = reminder.NextReminder.AddDate(0, 1, 0)
			} else {
				reminder.NextReminder = time.Date(2100, 1, 1, 0, 0, 0, 0, time.Local)
			}

			if err := database.DB.Save(&reminder).Error; err != nil {
				fmt.Printf("Ошибка при обновлении напоминания: %v\n", err)
			}
		}
	}
	return nil
}

func sendReminders(bot *tgbotapi.BotAPI, userID uint64, reminders []models.Reminder, log *slog.Logger) {
	chatID := int64(userID)

	if len(reminders) == 0 {
		return
	}

	text := "🛎 *Напоминание* \n\n Не забудьте сегодня совершить платеж(и):\n"
	for i, reminder := range reminders {
		text += fmt.Sprintf(
			"%d.%s\n  Сумма: %d\n\n",
			i+1,
			reminder.Category,
			reminder.Amount,
		)
	}
	text += "Хорошего дня! 🐙\n"

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	_, err := bot.Send(msg)
	if err != nil {
		log.Error(
			"Failed to send reminder to user %d:",
			"error", err)
	}
}
