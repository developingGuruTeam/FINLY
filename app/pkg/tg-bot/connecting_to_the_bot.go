package tg_bot

import (
	methods_for_user "cachManagerApp/app/internal/methods-for-user"
	"cachManagerApp/app/internal/notion"
	buttons_create "cachManagerApp/app/pkg/buttons-create"
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
			// базовый дефолтный старт!!!
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

			// сообщение от админов
			case "/send_admin_message":
				if isAdmin(update.Message.Chat.ID) { // Проверяем, является ли отправитель администратором
					message := "🎉 _Со Старым Новым годом!\nПусть этот год будет наполнен радостью, счастьем и финансовыми успехами!_ 🐙"
					SendOneTimeNotificationToAll(bot, message, log)
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🚫 У вас нет прав для выполнения этой команды.")
					bot.Send(msg)
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
