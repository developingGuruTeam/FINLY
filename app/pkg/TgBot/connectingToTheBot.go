package TgBot

import (
	"cachManagerApp/app/internal/methodsForUser"
	"cachManagerApp/app/internal/notion"
	"cachManagerApp/app/pkg/ButtonsCreate"
	"log/slog"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// подключение к тг и обработка обновлений
func ConnectToTgBot(log *slog.Logger) (*tgbotapi.BotAPI, error) {

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Error("Failed to connect to Telegram bot API: %v", err)
	}
	log.Info("Successfully connected to Telegram bot API!")

	bot.Debug = false

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 15

	updates := bot.GetUpdatesChan(updateConfig)

	// старт работы уведомлений
	notion.StartReminderServiceWithCron(bot, log)

	// старт всех кнопок
	buttonCreator := ButtonsCreate.TelegramButtonCreator{}

	for update := range updates {
		if update.Message != nil {
			switch update.Message.Text {
			case "/start":
				// планировщик отложенных уведомлений запускаем для пользователя вместе со стартом бота
				ScheduleNotifications(bot, update.Message.Chat.ID, update.Message.From.UserName, log)

				// высылаем только при старте /start
				mainMenuKeyboard := buttonCreator.CreateMainMenuButtons()
				userHandler := &methodsForUser.UserMethod{}
				if err := userHandler.PostUser(update, log); err != nil {
					log.Info("Ошибка при добавлении пользователя: %v", log.With("error", err))
				} else {
					log.Info("Пользователь успешно добавлен.", log.With("telegram_id", update.Message.Chat.ID))
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Добро пожаловать! 👋\nЯ — ваш финансовый помощник.\nБлагодаря мне у вас есть возможность взять свои денежные средства под контроль.\nВперёд к финансовому успеху!\nВыберите действие в меню ✏\nБазовые команды бота:\n/help - Помощь в использовании\n/hi - Мотивационное сообщение")
				msg.ReplyMarkup = mainMenuKeyboard
				if _, err := bot.Send(msg); err != nil {
					log.Error("Failed to send message with main menu buttons: %v", log.With("error", err))
				}
			default:
				// обработчик
				chatID := update.Message.Chat.ID
				if _, ok := notion.RemindersStates[chatID]; ok {
					// Если пользователь уже в процессе, обрабатываем его ввод
					notion.HandleReminderInput(bot, update, log)
				} else {
					PushOnButton(bot, update, buttonCreator, log)
				}
			}
		}
	}
	return bot, nil
}
