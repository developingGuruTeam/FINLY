package notion

import (
	"cachManagerApp/app/db/models"
	rules_for_notion "cachManagerApp/app/internal/notion/rules-for-notion"
	buttons_create "cachManagerApp/app/pkg/buttons-create"
	"cachManagerApp/database"
	"log/slog"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var RemindersStates = map[int64]*models.Reminder{}

func StartReminder(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID

	// создаем новый Reminder для пользователя
	RemindersStates[chatID] = &models.Reminder{
		UserID: uint64(chatID),
	}
}

func HandleReminderInput(bot *tgbotapi.BotAPI, update tgbotapi.Update, log *slog.Logger) {
	chatID := update.Message.Chat.ID
	reminder, exists := RemindersStates[chatID]
	if !exists {
		return
	}

	switch {
	case reminder.Frequency == "":
		// получаем частоту платежа
		if update.Message.Text == "⬅ В меню" {
			menuMain := buttons_create.TelegramButtonCreator{}
			back := menuMain.CreateMainMenuButtons()
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вы вернулись в главное меню")
			msg.ReplyMarkup = back
			if _, err := bot.Send(msg); err != nil {
				log.Error("Failed to send main menu: %v", "Error", err)
			}
			// Удаляем напоминание из состояния
			delete(RemindersStates, update.Message.Chat.ID)
			return
		}

		frequency := update.Message.Text
		if frequency != "7️⃣ Каждую неделю" && frequency != "🌙 Каждый месяц" {
			msg := tgbotapi.NewMessage(chatID, "Неверный ввод. Пожалуйста, выберите '7️⃣ Каждую неделю' или '🌙 Каждый месяц'.")
			_, _ = bot.Send(msg)
			return
		}

		if frequency == "7️⃣ Каждую неделю" {
			reminder.Frequency = "неделя"
		}

		if frequency == "🌙 Каждый месяц" {
			reminder.Frequency = "месяц"
		}

		// Переходим к следующему этапу — названию платежа
		msg := tgbotapi.NewMessage(chatID, "Введите название регулярного платежа\n_например: Кредит за машину_")
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true) // Убираем кнопки
		_, _ = bot.Send(msg)

	case reminder.Category == "":
		// Получаем название платежа
		reminder.Category = update.Message.Text

		msg := tgbotapi.NewMessage(chatID, "Введите дату следующего регулярного платежа (ДД.ММ.ГГГГ)\n_например: 01.02.2006_")
		msg.ParseMode = "Markdown"
		_, err := bot.Send(msg)
		if err != nil {
			log.Error("Ошибка в отправке сообщения в категории напоминания %v", "Error", err)
		}
		return

	case reminder.NextReminder.IsZero():
		// Проверяем и сохраняем дату платежа

		nextReminder, err := rules_for_notion.ValidateRightTime(update.Message.Text)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, err.Error())
			_, _ = bot.Send(msg)
			return
		}

		reminder.NextReminder = nextReminder

		msg := tgbotapi.NewMessage(chatID, "Введите сумму платежа")
		_, err = bot.Send(msg)
		if err != nil {
			log.Error("Ошибка в отправке сообщения суммы: %v", "Error", err)
		}

	case reminder.Amount == 0:
		// Получаем сумму платежа
		amount, err := strconv.Atoi(update.Message.Text)
		if err != nil || amount <= 0 {
			msg := tgbotapi.NewMessage(chatID, "Пожалуйста, введите целое положительное число.")
			_, _ = bot.Send(msg)
			return
		}

		reminder.Amount = amount
		reminder.CreatedAt = time.Now()

		if err := database.DB.Create(&reminder).Error; err != nil {
			log.Error("Ошибка при сохранении напоминания: %v", "Error", err)
			msg := tgbotapi.NewMessage(chatID, "Произошла ошибка при сохранении напоминания. Попробуйте позже.")
			_, _ = bot.Send(msg)
			return
		}

		menuMain := buttons_create.TelegramButtonCreator{}
		back := menuMain.CreateMainMenuButtons()
		msg := tgbotapi.NewMessage(chatID, "Напоминание успешно создано 🐙")

		msg.ReplyMarkup = back
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send main menu: %v", "Error", err)
		}

		// Удаляем напоминание из состояния
		delete(RemindersStates, chatID)
		return
	}
}
