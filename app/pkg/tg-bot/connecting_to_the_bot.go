package tg_bot

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
		log.Error("Failed to connect to Telegram bot API:", slog.Any("error", err))
	}
	log.Info("Successfully connected to Telegram bot API!")

	bot.Debug = false

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 15

	updates := bot.GetUpdatesChan(updateConfig)

	// старт работы уведомлений
	notion.StartReminderServiceWithCron(bot, log)

	// старт всех кнопок
	buttonCreator := buttons_create.TelegramButtonCreator{}

	for update := range updates {
		if update.Message != nil {
			switch update.Message.Text {
			case "/start":

				userHandler := &methods_for_user.UserMethod{}
				if err := userHandler.PostUser(update, log); err != nil {
					log.Info("Ошибка при добавлении пользователя:", slog.Any("error", err))
				} else {
					log.Info("Пользователь успешно добавлен.", slog.Any("user added", userHandler))
				}

				// планировщик отложенных уведомлений запускаем для пользователя вместе со стартом бота
				ScheduleNotifications(bot, update.Message.Chat.ID, update.Message.From.UserName, log)

				// отправляем стартовое сообщение
				WelcomeMessage(bot, update.Message.Chat.ID, buttonCreator, log)

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
