package TgBot

import (
	"cachManagerApp/app/internal/methodsForUser"
	"cachManagerApp/app/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
)

// подключение к тг и обработка обновлений
func ConnectToTgBot() (*tgbotapi.BotAPI, error) {
	log := logger.GetLogger()

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatalf("Failed to connect to Telegram bot API: %v", err)
	}
	log.Info("Successfully connected to Telegram bot API!")

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 15

	updates := bot.GetUpdatesChan(updateConfig)

	// старт всех кнопок
	buttonCreator := TelegramButtonCreator{}

	for update := range updates {
		if update.Message != nil {
			switch update.Message.Text {
			case "/start":
				// планировщик отложенных уведомлений запускаем для пользователя вместе со стартом бота
				ScheduleNotifications(bot, update.Message.Chat.ID, update.Message.From.UserName)

				// высылаем только при старте /start
				mainMenuKeyboard := buttonCreator.CreateMainMenuButtons()
				userHandler := &methodsForUser.UserMethod{}
				if err := userHandler.PostUser(update); err != nil {
					log.Printf("Ошибка при добавлении пользователя: %v", err)
				} else {
					log.Println("Пользователь успешно добавлен.")
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Добро пожаловать! 👋\nЯ — ваш финансовый помощник.\nБлагодаря мне у вас есть возможность взять свои денежные средства под контроль.\nВперёд к финансовому успеху!\nВыберите действие в меню ✏\n\nБазовые команды бота:\n/info - Информация о боте\n/help - Помощь в использовании\n/hi - Мотивационное сообщение")
				msg.ReplyMarkup = mainMenuKeyboard
				if _, err := bot.Send(msg); err != nil {
					log.Printf("Failed to send message with main menu buttons: %v", err)
				}
			default:
				// обработчик
				PushOnButton(bot, update, buttonCreator)
			}
		}
	}

	return bot, nil
}
