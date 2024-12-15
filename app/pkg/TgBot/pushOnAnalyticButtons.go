package TgBot

import (
	"cachManagerApp/app/internal/methodsForAnalytic/methodsForSummary"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
		report, err := methodsForSummary.AnalyseByCategoriesWeek(update)
		if err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось получить данные. Попробуйте позже.")
			_, _ = bot.Send(msg)
			log.Printf("Ошибка получения данных за день: %v", err)
			return
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, report)
		_, _ = bot.Send(msg)

	case "🤳 месяц":
		report, err := methodsForSummary.AnalyseByCategoriesMonth(update)
		if err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось получить данные. Попробуйте позже.")
			_, _ = bot.Send(msg)
			log.Printf("Ошибка получения данных за день: %v", err)
			return
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, report)
		_, _ = bot.Send(msg)
	}
}
