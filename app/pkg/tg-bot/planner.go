package tg_bot

import (
	"cachManagerApp/app/db/models"
	"cachManagerApp/database"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-co-op/gocron" // библиотека для работы с планировщиком
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// мапа для хранения активного планировщика
var activeSchedulers = make(map[int64]*gocron.Scheduler)

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
		log.Error("Ошибка отправки сообщения:", slog.Any("error", err))
	}
}

// Отправляем уведомления каждый день в 12:00 по среднеевропейскому времени
func ScheduleNotifications(bot *tgbotapi.BotAPI, chatID int64, userName string, log *slog.Logger) {
	// проверка запущен ли планировщик для пользователя
	if _, exists := activeSchedulers[chatID]; exists {
		log.Info("Уведомление уже запущено для пользователя", slog.Int64("chatID", chatID))
		return
	}
	scheduler := gocron.NewScheduler(time.Local)

	//scheduler.Cron("*/1 * * * *").Do(func() { // временная хрень для тестов уведомление раз в минуту НЕ УДАЛЯТЬ!
	scheduler.Cron("0 12 * * *").Do(func() { // БОЕВАЯ ЧАСТЬ! это уведомление на каждый день в 12 (по факту 15 часа для мск)
		SendNotificationToUser(bot, chatID, userName, log) // Отправляем уведомление
	})

	// асинхронный запуск планировщика!
	scheduler.StartAsync()
	activeSchedulers[chatID] = scheduler
}

// разовое уведомление всем пользователям для админ команды
func SendOneTimeNotificationToAll(bot *tgbotapi.BotAPI, message string, log *slog.Logger) {
	// извлекаем всех пользователей из БД
	var users []models.Users
	result := database.DB.Find(&users)
	if result.Error != nil {
		log.Error("Ошибка получения списка пользователей из базы данных", slog.Any("error", result.Error))
		return
	}

	// отправляем уведомление каждому пользователю в цикле
	for _, user := range users {
		chatID := int64(user.TelegramID) // приводим TelegramID к int64
		msg := tgbotapi.NewMessage(chatID, message)
		msg.ParseMode = "Markdown"
		_, err := bot.Send(msg)
		if err != nil {
			log.Error("Ошибка отправки разового уведомления пользователю", slog.Int64("chatID", chatID))
		} else {
			log.Info("Разовое уведомление отправлено пользователю", slog.Int64("chatID", chatID))
		}
	}
}

// проверка на админку
func isAdmin(chatID int64) bool {
	admins := os.Getenv("ADMIN_CHAT_IDS")  // список администраторов из .env
	adminIDs := strings.Split(admins, ",") // разбираем строку на отдельные ID

	for _, id := range adminIDs {
		userID, err := strconv.ParseInt(strings.TrimSpace(id), 10, 64) // преобразуем строку в int64
		// проверяем, совпадает ли юзер с одним из админов
		if err == nil && userID == chatID {
			return true
		}
	}
	return false
}
