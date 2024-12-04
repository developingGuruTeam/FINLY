package TgBot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math/rand"
)

// интерфейс создания кнопок
type ButtonCreator interface {
	CreateMainMenuButtons() tgbotapi.ReplyKeyboardMarkup
	//CreateInlineButtons() tgbotapi.InlineKeyboardMarkup
}

// структура интерфейса создания кнопок
type TelegramButtonCreator struct{}

// cоздание кнопок главного меню по строкам
func (t TelegramButtonCreator) CreateMainMenuButtons() tgbotapi.ReplyKeyboardMarkup {
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

func (t TelegramButtonCreator) CreateIncomeMenuButtons() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📥 Заработная плата"),
			tgbotapi.NewKeyboardButton("📤 Дополнительный доход"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📥 Доход от бизнеса"),
			tgbotapi.NewKeyboardButton("📤 Доход от инвестиций"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📥 Гос выплаты"),
			tgbotapi.NewKeyboardButton("📤 Продажа имущества"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📥 Прочее"),
			tgbotapi.NewKeyboardButton("⬅ В меню"),
		),
	)
}

// cоздание кнопок меню отчета по строкам
func (t TelegramButtonCreator) CreateReportsMenuButtons() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📑 Отчет за день"),
			tgbotapi.NewKeyboardButton("📑 Отчет за неделю"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📑 Отчет за месяц"),
			tgbotapi.NewKeyboardButton("⬅ В меню"),
		),
	)
}

// cоздание кнопок меню настроек по строкам
func (t TelegramButtonCreator) CreateSettingsMenuButtons() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🎭 Изменить имя"),
			tgbotapi.NewKeyboardButton("💫 Тарифный план"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("💱 Изменить валюту"),
			tgbotapi.NewKeyboardButton("⬅ В меню"),
		),
	)
}

//// создание inline кнопок
//func (t TelegramButtonCreator) CreateInlineButtons() tgbotapi.InlineKeyboardMarkup {
//	return tgbotapi.NewInlineKeyboardMarkup(
//		tgbotapi.NewInlineKeyboardRow(
//			tgbotapi.NewInlineKeyboardButtonData("info", "info"),
//			tgbotapi.NewInlineKeyboardButtonData("help", "help"),
//			tgbotapi.NewInlineKeyboardButtonData("hi", "hi"),
//		),
//	)
//}

// рандомное сообщение для команды /hi
func randomTextForHi() string {
	hiText := [...]string{
		"💰 Сегодня отличный день, чтобы начать экономить!", "💰 Ну что, приступим считать твои траты?",
		"💰 Интересно, сегодня день будет ➕ или ➖ ?", "💰 Экономия должна быть экономной!",
		"💰 Сэкономил, значит заработал!", "💰 Время взять финансы под контроль!",
		"💰 Денежки любят счёт. Начнём?", "💰 Финансовый учёт — первый шаг к успеху!",
		"💰 Сегодня день больших возможностей для экономии!", "💰 Планируй расходы — станешь богаче!",
		"💰 Контроль за тратами — твоя суперсила!", "💰 Давайте посмотрим, куда уходят ваши денежки!",
		"💰 Финансовый контроль — это просто. Поехали!", "💰 Посчитаем твои финансы и найдём резервы!",
		"💰 Чем меньше траты, тем больше возможностей!", "💰 Успех начинается с грамотного управления финансами!",
		"💰 Настройся на экономию и достигай целей!", "💰 Каждый шаг к учёту — шаг к финансовой свободе!",
		"💰 Добро пожаловать в мир финансового порядка!", "💰 Везде нужен порядок. Особенно в деньгах!",
	}
	return hiText[rand.Intn(len(hiText))]
}
