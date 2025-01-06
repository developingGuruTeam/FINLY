package TgBot

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/go-co-op/gocron" // библиотека для работы с планировщиком
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Отправляем уведомление пользователю
func SendNotificationToUser(bot *tgbotapi.BotAPI, chatID int64, userName string, log *slog.Logger) {
	clearName, _ := ClearUserNameFromChatID(chatID)
	if clearName == "" {
		clearName = userName
	}
	message := fmt.Sprintf("Привет, %s 👋!\nНе забывай записывать свои приходы и расходы, чтобы вести их учет 🧮", clearName)

	// Создаем объект сообщения
	msg := tgbotapi.NewMessage(chatID, message)

	// Отправляем сообщение пользователю
	if _, err := bot.Send(msg); err != nil {
		log.Error("Ошибка отправки сообщения: %v", err)
	}
}

// Отправляем уведомления каждый день в 12:00 по среднеевропейскому времени
func ScheduleNotifications(bot *tgbotapi.BotAPI, chatID int64, userName string, log *slog.Logger) {
	// Создаем новый планировщик
	scheduler := gocron.NewScheduler(time.Local) // use local time without timezone

	// scheduler.Cron("*/1 * * * *").Do(func() { // временная хрень для тестов уведомление раз в минуту НЕ УДАЛЯТЬ!
	scheduler.Cron("0 12 * * *").Do(func() { // БОЕВАЯ ЧАСТЬ! это уведомление на каждый день в 12
		SendNotificationToUser(bot, chatID, userName, log) // Отправляем уведомление
	})

	// асинхронный запуск планировщика!
	scheduler.StartAsync()
}
