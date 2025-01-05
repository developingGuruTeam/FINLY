package methodsForUser

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
		log.Info("Пользователь существует", log.With("telegram_id", user.TelegramID))
		return errors.New("user already exists")
	}

	if err := database.DB.Create(&user).Error; err != nil {
		log.Error("Ошибка добавления нового пользователя: %v", log.With("Error", err))
		return err
	}

	log.Info("Новый пользователь успешно добавлен.", log.With("telegram_id", user.TelegramID))
	return nil
}

func (u *UserMethod) UpdateUserName(update tgbotapi.Update) error {
	newUserName := update.Message.Text

	res := database.DB.Model(&models.Users{}).
		Where("telegram_id = ?", uint64(update.Message.Chat.ID)).
		Update("name", newUserName)
	if res.Error != nil {
		log.Printf("Ошибка обновления имени пользователя: %v", res.Error)
		return res.Error
	}

	// Проверяем, обновлена ли хотя бы одна строка
	if res.RowsAffected == 0 {
		log.Println("Не найден пользователь с указанным telegram_id.")
		return errors.New("пользователь не найден")
	}

	log.Printf("Имя пользователя с Telegram ID %d успешно обновлено на '%s'.", update.Message.Chat.ID, newUserName)

	return nil
}

func (u *UserMethod) UpdateUserCurrency(update tgbotapi.Update) error {
	newCurrency := strings.ToLower(update.Message.Text)

	// делаем тут валюту из трех букв и точку на конце
	runes := []rune(newCurrency)
	if len(runes) > 3 {
		runes = runes[:3]
	}
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
