package methods_for_user

import (
	"cachManagerApp/app/db/models"
	"cachManagerApp/database"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//go:generate mockery --name=UsersHandlers --output=../tests/mocks --with-expecter
type UsersHandlers interface {
	PostUser(update tgbotapi.Update) error
	UpdateUserName(update tgbotapi.Update) error
	UpdateUserCurrency(update tgbotapi.Update) error
}

type UserMethod struct{}

func (u *UserMethod) PostUser(update tgbotapi.Update, log *slog.Logger) error {
	user := models.Users{
		TelegramID: uint64(update.Message.Chat.ID),
		Name:       update.Message.From.UserName,
		Currency:   "руб.",
	}

	var userExist models.Users
	res := database.DB.Where("telegram_id = ?", user.TelegramID).First(&userExist)

	if res.Error == nil {
		log.Info("Пользователь существует", "telegram_id", user.TelegramID)
		return errors.New("user already exists")
	}

	if err := database.DB.Create(&user).Error; err != nil {
		log.Error("Ошибка добавления нового пользователя: %v", "Error", err)
		return err
	}

	log.Info("Новый пользователь успешно добавлен.", "telegram_id", user.TelegramID)
	return nil
}

func (u *UserMethod) UpdateUserName(update tgbotapi.Update) error {
	newUserName := update.Message.Text
	telegramID := uint64(update.Message.Chat.ID)

	// проверяем текущее имя
	var currentUser models.Users
	res := database.DB.First(&currentUser, "telegram_id = ?", telegramID)
	if res.Error != nil {
		log.Printf("Ошибка поиска пользователя: %v", res.Error)
		return res.Error
	}

	// обновляем имя
	res = database.DB.Model(&models.Users{}).
		Where("telegram_id = ?", telegramID).
		Update("name", newUserName)
	if res.Error != nil {
		log.Printf("Ошибка обновления имени пользователя: %v", res.Error)
		return res.Error
	}

	if res.RowsAffected == 0 {
		log.Println("Не найден пользователь с указанным telegram_id.")
		return errors.New("пользователь не найден")
	}

	log.Printf("Имя пользователя с Telegram ID %d успешно обновлено на '%s'.", telegramID, newUserName)
	return nil
}

func (u *UserMethod) UpdateUserCurrency(update tgbotapi.Update) error {
	newCurrency := strings.ToLower(update.Message.Text)

	// делаем тут валюту из трех букв и точку на конце
	runes := []rune(newCurrency)
	newCurrency = fmt.Sprintf("%s.", string(runes))

	res := database.DB.Model(&models.Users{}).
		Where("telegram_id = ?", uint64(update.Message.Chat.ID)).
		Update("currency", newCurrency)

	if res.Error != nil {
		log.Printf("Ошибка обновления валюты: %v", res.Error)
		return res.Error
	}

	if res.RowsAffected == 0 {
		log.Println("Не найден пользователь с указанным telegram_id.")
		return errors.New("пользователь не найден")
	}

	log.Printf("Валюта успешна обновлена на '%s'.", newCurrency)
	return nil
}
