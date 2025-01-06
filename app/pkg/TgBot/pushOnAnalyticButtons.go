package TgBot

import (
	"cachManagerApp/app/internal/methodsForAnalytic/methodsForSummary"
	"cachManagerApp/app/internal/notion"
	"cachManagerApp/app/pkg/ButtonsCreate"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func PushOnAnalyticButton(bot *tgbotapi.BotAPI, update tgbotapi.Update, buttonCreator ButtonsCreate.TelegramButtonCreator, command string, log *slog.Logger) {
	currency, _ := CurrencyFromChatID(update.Message.Chat.ID)

	switch command {

	case "⚖️ Cальдо":
		saldo := buttonCreator.CreateSaldoAnalyticButtons()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите период")
		msg.ReplyMarkup = saldo
		if _, err := bot.Send(msg); err != nil {
			log.Info("Failed to send main menu: %v", log.With("error", err))
		}

	case "💲Сальдо за неделю":
		summary, err := methodsForSummary.AnalyseBySaldoWeek(update)
		if err != nil {
			log.Info("Failed to get summary in the week period: %v", log.With("error", err))
		}
		response := methodsForSummary.GenerateWeeklySaldoReport(summary, currency)
		newMsg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
		newMsg.ParseMode = "Markdown"
		_, _ = bot.Send(newMsg)

	case "💰Сальдо за месяц":
		summary, err := methodsForSummary.AnalyseBySaldoMonth(update)
		if err != nil {
			log.Info("Failed to get summary in the month period: %v", log.With("error", err))
		}
		response := methodsForSummary.GenerateMonthlySaldoReport(summary, currency)
		newMsg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
		newMsg.ParseMode = "Markdown"
		_, _ = bot.Send(newMsg)

	// напоминания об оплате
	case "💡 Напоминание":
		notion := buttonCreator.CreateNotionButtons()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите тип напоминания")
		msg.ReplyMarkup = notion
		if _, err := bot.Send(msg); err != nil {
			log.Info("Failed to send main menu: %v", log.With("error", err))
		}

	case "📅 Регулярный платёж":
		// создаем мапу для работы с напоминаниями
		notion.StartReminder(bot, update)
		reminder := buttonCreator.CreateFreqButtons()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите периодичность платежа:")
		msg.ReplyMarkup = reminder
		_, err := bot.Send(msg)
		if err != nil {
			log.Error("Error sending message: %v", err)
		}
	}
}
