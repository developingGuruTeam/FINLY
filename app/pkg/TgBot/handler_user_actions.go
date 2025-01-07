package TgBot

import (
	"cachManagerApp/app/db/models"
	"cachManagerApp/app/internal/methodsForUser"
	"cachManagerApp/app/pkg/ButtonsCreate"
	"cachManagerApp/database"
	"fmt"
	"log/slog"
	"strings"
	"unicode"
	"unicode/utf8"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// обработчик действий пользователя для изменения имени и валюты
func handleUserAction(bot *tgbotapi.BotAPI, update tgbotapi.Update, userResp UserResponse, buttonCreator ButtonsCreate.TelegramButtonCreator, log *slog.Logger) {
	chatID := update.Message.Chat.ID

	switch userResp.Action {

	case "change_name":
		newName := strings.TrimSpace(update.Message.Text) // убираем пробелы по обе стороны, если есть
		if newName == "" {
			newName = "Пользователь"
		}

		// проверка нового имени: только буквы и длина от 1 до 32 символов
		var validName bool = true
		for _, symbol := range newName {
			if !unicode.IsLetter(symbol) && symbol != ' ' { // имя только из букв и пробелов
				validName = false
				break
			}
		}

		if utf8.RuneCountInString(newName) == 0 || utf8.RuneCountInString(newName) > 32 || !validName {
			msg := tgbotapi.NewMessage(chatID, "🚫 Некорректное имя. Имя должно содержать только буквы и быть не более 32 символов.")
			bot.Send(msg)
			return
		}

		// обновляем имя пользователя
		user := methodsForUser.UserMethod{}
		if err := user.UpdateUserName(update); err != nil {
			log.Error("Ошибка обновления имени пользователя", log.With("error", err))
			msg := tgbotapi.NewMessage(chatID, "❌ Ошибка при обновлении имени.")
			bot.Send(msg)
			return
		}

		msgDone := fmt.Sprintf("✅ Ваше имя успешно изменено на *%s*.", newName)
		returnToMainMenu(bot, chatID, buttonCreator, msgDone)

	case "change_currency":
		newCurrency := strings.ToLower(update.Message.Text) // преобразуем в нижний регистр
		// проверка новой валюты на алфавит
		var okCurrency bool = true
		for _, symbol := range newCurrency {
			if !unicode.IsLetter(symbol) {
				okCurrency = false
				break
			}
		}
		// проверка валюты на длину
		if utf8.RuneCountInString(newCurrency) != 3 || okCurrency != true {
			msg := tgbotapi.NewMessage(chatID, "🚫 Некорректный формат валюты. Валюта должна содержать только буквы и быть не более 3 символов.")
			bot.Send(msg)
			return
		}
		// обновляем валюту пользователя
		user := methodsForUser.UserMethod{}
		if err := user.UpdateUserCurrency(update); err != nil {
			log.Error("Ошибка обновления валюты", log.With("error", err))
			msg := tgbotapi.NewMessage(chatID, "❌ Ошибка при обновлении валюты")
			bot.Send(msg)
			return
		}

		msgDone := fmt.Sprintf("✅ Ваша валюта изменена на *%s*.", newCurrency)
		returnToMainMenu(bot, chatID, buttonCreator, msgDone)
	}

	mu.Lock()
	delete(userStates, chatID) // удаляем состояние после обработки
	mu.Unlock()
}

// возврат кнопок меню и удаления состояния после обработки транзакции
func returnToMainMenu(bot *tgbotapi.BotAPI, chatID int64, buttonCreator ButtonsCreate.TelegramButtonCreator, msg string) {
	// создаем кнопки главного меню
	mainMenu := buttonCreator.CreateMainMenuButtons()

	// отправляем пустое сообщение с кнопками
	menuMsg := tgbotapi.NewMessage(chatID, msg)
	menuMsg.ParseMode = "Markdown"
	menuMsg.ReplyMarkup = mainMenu // показываем кнопки
	bot.Send(menuMsg)

	delete(transactionStates, chatID)
}

// получение имени из БД
func ClearUserNameFromChatID(chatID int64) (string, error) {
	var user models.Users
	result := database.DB.Where("telegram_id = ?", chatID).First(&user)
	if result.Error != nil {
		return "", result.Error
	}
	fmt.Println(user.Name)
	return user.Name, nil
}

// получение валюты из бд
func CurrencyFromChatID(chatID int64) (string, error) {
	var user models.Users
	result := database.DB.Where("telegram_id = ?", chatID).First(&user)
	if result.Error != nil {
		return "", result.Error
	}
	return user.Currency, nil
}
