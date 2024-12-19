package TgBot

import (
	"cachManagerApp/app/db/models"
	"cachManagerApp/app/internal/methodsForAnalytic/methodsForSummary"
	"cachManagerApp/database"
	redisDB "cachManagerApp/database/redis"
	"context"
	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

func PushOnAnalyticButton(bot *tgbotapi.BotAPI, update tgbotapi.Update, buttonCreator TelegramButtonCreator, command string) {
	currency, _ := CurrencyFromChatID(update.Message.Chat.ID)

	switch command {
	case "🛍 Анализ категорий":
		category := buttonCreator.CreateCategoryAnalyticButtons()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите период")
		msg.ReplyMarkup = category
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send main menu: %v", err)
		}

	case "💲Анализ за неделю":
		redisClient, err := redisDB.NewRedisClient()
		if err != nil {
			log.Infof("Failed to connect to Redis: %v", err)
		}
		key := update.Message.Text + update.Message.Chat.UserName
		report, err := redisClient.Client.Get(context.Background(), key).Result()
		if err == redis.Nil {
			report, err = methodsForSummary.AnalyseByCategoriesWeek(update)
			redisClient.Client.Set(context.Background(), key, report, time.Hour)
			time.Sleep(2 * time.Second)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось получить данные. Попробуйте позже.")
				_, _ = bot.Send(msg)
				log.Printf("Ошибка получения данных за день: %v", err)
				return
			}

		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, report)
		_, _ = bot.Send(msg)

	case "💰Анализ за месяц":
		key := update.Message.Text + update.Message.Chat.UserName
		redisClient, err := redisDB.NewRedisClient()
		if err != nil {
			log.Infof("Failed to connect to Redis: %v", err)
		}
		report, err := redisClient.Client.Get(context.Background(), key).Result()
		if err == redis.Nil {
			report, err = methodsForSummary.AnalyseByCategoriesMonth(update)
			redisClient.Client.Set(context.Background(), key, report, time.Hour)
			time.Sleep(2 * time.Second)
			log.Println("wait))")
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось получить данные. Попробуйте позже.")
				_, _ = bot.Send(msg)
				log.Printf("Ошибка получения данных за день: %v", err)
				return
			}
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, report)
		_, _ = bot.Send(msg)

	case "сальдо":
		saldo := buttonCreator.CreateSaldoAnalyticButtons()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите период")
		msg.ReplyMarkup = saldo
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send main menu: %v", err)
		}

	case "💲Сальдо за неделю":
		summary, err := methodsForSummary.AnalyseBySaldoWeek(update)
		if err != nil {
			log.Printf("Failed to get summary in the week period: %v", err)
		}
		response := methodsForSummary.GenerateWeeklySaldoReport(summary, currency)
		newMsg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
		newMsg.ParseMode = "Markdown"
		_, _ = bot.Send(newMsg)

	case "💰Сальдо за месяц":
		summary, err := methodsForSummary.AnalyseBySaldoMonth(update)
		if err != nil {
			log.Printf("Failed to get summary in the month period: %v", err)
		}
		response := methodsForSummary.GenerateMonthlySaldoReport(summary, currency)
		newMsg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
		newMsg.ParseMode = "Markdown"
		_, _ = bot.Send(newMsg)
	}
}

// Получение валюты из бд
func CurrencyFromChatID(chatID int64) (string, error) {
	var user models.Users
	result := database.DB.Where("telegram_id = ?", chatID).First(&user)
	if result.Error != nil {
		return "", result.Error
	}
	return user.Currency, nil
}
