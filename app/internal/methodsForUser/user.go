package methodsForUser

import (
	"cachManagerApp/app/db/models"
	"cachManagerApp/database"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type UsersHandlers interface {
	PostUser(update tgbotapi.Update) error
}

type UserMethod struct{}

func (u *UserMethod) PostUser(update tgbotapi.Update) error {
	user := models.Users{
		TelegramID: uint64(update.Message.Chat.ID),
		Name:       update.Message.From.UserName,
		Currency:   "RUB",
	}

	var userExist models.Users
	res := database.DB.Where("telegram_id = ?", user.TelegramID).First(&userExist)

	if res.Error == nil {
		log.Println("Пользователь существует")
		return errors.New("user already exists")
	}

	if err := database.DB.Create(&user).Error; err != nil {
		log.Printf("Ошибка добавления нового пользователя: %v", err)
		return err
	}

	log.Println("Новый пользователь успешно добавлен.")
	return nil
}
