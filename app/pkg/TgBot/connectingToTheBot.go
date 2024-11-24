package TgBot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

func ConnectToTgBot() (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatalf("Failed to connect to Telegram bot API: %v", err)
	}

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 15

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "hello world")

		switch update.Message.Command() {
		case "help":
			msg.Text = "..."
		case "hi":
			msg.Text = "Даров :)"
		case "start":
			msg.Text = "Я ещё не совсем готов, но можно потестить меню"

		case "bye":
			msg.Text = "Давай делай падла!!"
		default:
			msg.Text = "I don't know that command"
		}

		if _, err := bot.Send(msg); err != nil {

			log.Fatalf("Failed to send message: %v", err)
		}
	}

	return bot, nil
}
