package tg_bot

import (
	methods_for_summary "cachManagerApp/app/internal/methods-for-analytic/methods-for-summary"
	"cachManagerApp/app/internal/notion"
	buttons_create "cachManagerApp/app/pkg/buttons-create"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func PushOnAnalyticButton(bot *tgbotapi.BotAPI, update tgbotapi.Update, buttonCreator buttons_create.TelegramButtonCreator, command string, log *slog.Logger) {
	currency, _ := CurrencyFromChatID(update.Message.Chat.ID)

	switch command {

	case "⚖️ Cальдо":
		saldo := buttonCreator.CreateSaldoAnalyticButtons()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите период")
		msg.ReplyMarkup = saldo
		if _, err := bot.Send(msg); err != nil {
			log.Info("Failed to send main menu: %v", slog.Any("error", err))
		}

	case "💲Сальдо за неделю":
		summary, err := methods_for_summary.AnalyseBySaldoWeek(update)
		if err != nil {
			log.Info("Failed to get summary in the week period: %v", slog.Any("error", err))
		}
		response := methods_for_summary.GenerateWeeklySaldoReport(summary, currency)
		newMsg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
		newMsg.ParseMode = "Markdown"
		_, _ = bot.Send(newMsg)

	case "💰Сальдо за месяц":
		summary, err := methods_for_summary.AnalyseBySaldoMonth(update)
		if err != nil {
			log.Info("Failed to get summary in the month period: %v", slog.Any("error", err))
		}
		response := methods_for_summary.GenerateMonthlySaldoReport(summary, currency)
		newMsg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
		newMsg.ParseMode = "Markdown"
		_, _ = bot.Send(newMsg)

	// напоминания об оплате
	case "🛎 Напоминание":
		notion := buttonCreator.CreateNotionButtons()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите тип напоминания")
		msg.ReplyMarkup = notion
		if _, err := bot.Send(msg); err != nil {
			log.Info("Failed to send main menu: %v", slog.Any("error", err))
		}

	case "🔁 Регулярный платёж":
		// создаем мапу для работы с напоминаниями
		notion.StartReminder(bot, update)
		reminder := buttonCreator.CreateFreqButtons()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите периодичность платежа:")
		msg.ReplyMarkup = reminder
		_, err := bot.Send(msg)
		if err != nil {
			log.Error("Error sending message:", slog.Any("error", err))
		}
	}
}
