package TgBot

import (
	"cachManagerApp/app/internal/methodsForAnalytic/methodsForSummary"
	redisDB "cachManagerApp/database/redis"
	"context"
	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

func PushOnAnalyticButton(bot *tgbotapi.BotAPI, update tgbotapi.Update, buttonCreator TelegramButtonCreator, command string) {
	switch command {
	case "🛍 По категориям":
		category := buttonCreator.CreateCategoryAnalyticButtons()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите период")
		msg.ReplyMarkup = category
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send main menu: %v", err)
		}

	case "💅 неделя":
		redisClient, err := redisDB.NewRedisClient()
		if err != nil {
			log.Info("Failed to connect to Redis: %v", err)
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

	case "🤳 месяц":
		key := update.Message.Text + update.Message.Chat.UserName
		redisClient, err := redisDB.NewRedisClient()
		if err != nil {
			log.Info("Failed to connect to Redis: %v", err)
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
	}
}
