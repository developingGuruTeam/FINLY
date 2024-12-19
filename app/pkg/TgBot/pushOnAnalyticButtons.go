package TgBot

import (
	"cachManagerApp/app/db/models"
	"cachManagerApp/app/internal/methodsForAnalytic/methodsForSummary"
	"cachManagerApp/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func PushOnAnalyticButton(bot *tgbotapi.BotAPI, update tgbotapi.Update, buttonCreator TelegramButtonCreator, command string) {
	currency, _ := CurrencyFromChatID(update.Message.Chat.ID)

	switch command {
	case "💡 Напоминание":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "👷‍🔧`В разработке ...`\n\n`Ожидаемая дата выхода обновления 20.01.2025`")
		msg.ParseMode = "Markdown"
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send main menu: %v", err)
		}

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
