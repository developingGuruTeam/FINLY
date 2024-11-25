package TgBot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

// подключение к тг и обработка обновлений
func ConnectToTgBot() (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatalf("Failed to connect to Telegram bot API: %v", err)
	}

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 15

	updates := bot.GetUpdatesChan(updateConfig)

	// экземпляр TelegramStaticButtonCreator
	buttonCreator := TelegramStaticButtonCreator{}

	for update := range updates {
		if update.Message != nil {
			// выводит сообщение для статичного меню
			mainMenuKeyboard := buttonCreator.CreateMainMenuButtons()
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите действие")
			msg.ReplyMarkup = mainMenuKeyboard
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Failed to send message with main menu buttons: %v", err)
			}

			// для того чтобы активировать inline кнопки со слэшем
			inlineKeyboard := buttonCreator.CreateInlineInfoHelpButtons()
			msgInline := tgbotapi.NewMessage(update.Message.Chat.ID, "") // затер сообщение для обработки инлайн кнопок?, здесь по идее должно отображаться баланс на данный момент
			msgInline.ReplyMarkup = inlineKeyboard
			if _, err := bot.Send(msgInline); err != nil {
				log.Printf("Failed to send message with inline buttons: %v", err)
			}
		}

		// Обрабатываем нажатие на кнопки
		PushOnButton(bot, update)
	}

	return bot, nil
}
