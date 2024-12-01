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
	UpdateUserName(update tgbotapi.Update) error
}

type UserMethod struct {
	WaitingUpdate bool
}

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

func (u *UserMethod) UpdateUserName(update tgbotapi.Update) error {
	if u.WaitingUpdate {
		newUserName := update.Message.Text
		res := database.DB.Model(&models.Users{}).
			Where("telegram_id = ?", uint64(update.Message.Chat.ID)).
			Update("name", newUserName)
		if res.Error != nil {
			log.Printf("Ошибка обновления имени пользователя: %v", res.Error)
			return res.Error
		}
		log.Println("Имя пользователя успешно обновлено.")
		u.WaitingUpdate = false // Сбрасываем флаг ожидания
		return nil
	}
	return errors.New("no user name update expected")
}
