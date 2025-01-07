package tg_bot

import (
	"cachManagerApp/app/pkg/ButtonsCreate"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Информация о боте
func AboutBot(bot *tgbotapi.BotAPI, chatID int64, log *slog.Logger) {
	infoMessage := `
🐙 *FINLY - твой помощник в мире финансов* 

▪ Учет доходов и расходов 
▪ Информативные отчеты с графиками
▪ Напоминания о предстоящих платежах 

Все это позволит тебе взять финансы под свой контроль!

*Как пользоваться?*

1️⃣ Нажмите на кнопку *📥 Приход* для записи доходов.
2️⃣ Нажмите на кнопку *📤 Расход* для учета расходов.
3️⃣ Информацию о состоянии баланса можно получить в меню *📊 Отчеты*.
4️⃣ Управляйте настройками через меню *⚙️ Настройки*.
5️⃣ Получайте советы или мотивационные сообщения с помощью команды */hi*.

_Дополнительная информация:_
▪ При добавлении транзакции вы можете оставить комментарий или пропустить этот шаг.
▪ В меню *🛎 Напоминание* имеется способ настроить уведомление, чтобы не пропустить важный платеж.
▪ Все данные надежно сохраняются и доступны в удобной форме.


_© 2024-2025 г. Все права защищены._
`

	msg := tgbotapi.NewMessage(chatID, infoMessage)
	msg.ParseMode = "Markdown"

	// Отправляем сообщение пользователю
	if _, err := bot.Send(msg); err != nil {
		log.Info("Failed to send about bot message: ", slog.Any("error", err))
	}
}

func WelcomeMessage(bot *tgbotapi.BotAPI, chatID int64, buttonCreator buttons_create.TelegramButtonCreator, log *slog.Logger) {
	// создаем главное меню
	mainMenuKeyboard := buttonCreator.CreateMainMenuButtons()

	// сообщение приветствия
	welcomeMessage := "Добро пожаловать в *FINLY 🐙*!\n\nНажмите *ℹ️ Информация* в меню, чтобы узнать больше о моих возможностях.\n\nВперёд к финансовому успеху 🚀"
	msg := tgbotapi.NewMessage(chatID, welcomeMessage)
	msg.ParseMode = "Markdown" // включаем поддержку Markdown
	msg.ReplyMarkup = mainMenuKeyboard

	// отправляем сообщение
	if _, err := bot.Send(msg); err != nil {
		log.Error("Failed to send message with main menu buttons:", slog.Any("error", err))
	}
}
