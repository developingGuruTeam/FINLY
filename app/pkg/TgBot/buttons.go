package TgBot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// интерфейс создания кнопок
type StaticButtonCreator interface {
	CreateMainMenuButtons() tgbotapi.ReplyKeyboardMarkup
	CreateInlineInfoHelpButtons() tgbotapi.InlineKeyboardMarkup
}

// структура интерфейса создания кнопок
type TelegramStaticButtonCreator struct{}

// cоздание кнопок меню по строкам
func (t TelegramStaticButtonCreator) CreateMainMenuButtons() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📥 Приход"),
			tgbotapi.NewKeyboardButton("📤 Расход"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📊 Отчеты"),
			tgbotapi.NewKeyboardButton("⚙ Настройки"),
		),
	)
}

// создание inline кнопок через слэш
func (t TelegramStaticButtonCreator) CreateInlineInfoHelpButtons() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Info", "info"),
			tgbotapi.NewInlineKeyboardButtonData("Help", "help"),
		),
	)
}
