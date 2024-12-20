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
			tgbotapi.NewKeyboardButton("⚙️ Настройки"),
		),
	)
}

func (t TelegramButtonCreator) CreateIncomeMenuButtons() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("💳 Заработная плата"),
			tgbotapi.NewKeyboardButton("💱 Побочный доход"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("😎 Доход от бизнеса"),
			tgbotapi.NewKeyboardButton("🏦 Доход от инвестиций"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("👮‍ Гос. выплаты"),
			tgbotapi.NewKeyboardButton("🏠 Продажа имущества"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("⚪️ Прочие доходы"),
			tgbotapi.NewKeyboardButton("⬅ В меню"),
		),
	)
}

func (t TelegramButtonCreator) CreateExpensesMenuButtons() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🛍 Бытовые траты"),
			tgbotapi.NewKeyboardButton("🫡 Регулярные платежи"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("👘 Одежда"),
			tgbotapi.NewKeyboardButton("💪 Здоровье"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("👨‍🏫 Досуг и образование"),
			tgbotapi.NewKeyboardButton("🚀 Инвестиции"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("⚪️ Прочие расходы"),
			tgbotapi.NewKeyboardButton("⬅ В меню"),
		),
	)
}

// cоздание кнопок меню отчета по строкам
func (t TelegramButtonCreator) CreateReportsMenuButtons() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("💵 Отчет по доходам"),
			tgbotapi.NewKeyboardButton("💸 Отчет по расходам"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🧑‍💻 Аналитика"),
			tgbotapi.NewKeyboardButton("⬅ В меню"),
		),
	)
}

func (t TelegramButtonCreator) CreateIncomeAnalyticButtons() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📈 Отчет за день"),
			tgbotapi.NewKeyboardButton("📈 Отчет за неделю"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📈 Отчет за месяц"),
			tgbotapi.NewKeyboardButton("⬅ В меню"),
		),
	)
}

func (t TelegramButtonCreator) CreateExpensesAnalyticButtons() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📉 Отчет за день"),
			tgbotapi.NewKeyboardButton("📉 Отчет за неделю"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📉 Отчет за месяц"),
			tgbotapi.NewKeyboardButton("⬅ В меню"),
		),
	)
}

func (t TelegramButtonCreator) CreateSuperAnalyticButtons() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🧮 Статистика"),
			tgbotapi.NewKeyboardButton("🤑 Cальдо"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("👨‍🔬 Экспертная аналитика"),
			tgbotapi.NewKeyboardButton("⬅ В меню"),
		),
	)
}

func (t TelegramButtonCreator) CreateSaldoAnalyticButtons() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("💲Сальдо за неделю"),
			tgbotapi.NewKeyboardButton("💰Сальдо за месяц"),
		),
		tgbotapi.NewKeyboardButtonRow(
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
			tgbotapi.NewKeyboardButton("💡 Создать напоминание"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("⬅ В меню"),
		),
	)
}

// кнопки с уведомлениями
func (t TelegramButtonCreator) CreateNotionButtons() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📅 Регулярный платёж"),
			tgbotapi.NewKeyboardButton("🎯 Накопления"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🛒 Одно напоминание"),
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
		"💰 Интересно, сегодня день будет ➕ или ➖ ?", "💰 Экономика должна быть экономной!",
		"💰 Сэкономил, значит заработал!", "💰 Время взять финансы под контроль!",
		"💰 Денежки любят счёт. Начнём?", "💰 Финансовый учёт — первый шаг к успеху!",
		"💰 Сегодня день больших возможностей для экономии!", "💰 Планируй расходы — станешь богаче!",
		"💰 Контроль за тратами — твоя суперсила!", "💰 Давайте посмотрим, куда уходят ваши денежки!",
		"💰 Финансовый контроль — это просто. Поехали!", "💰 Посчитаем твои финансы и найдём резервы!",
		"💰 Чем меньше траты, тем больше возможностей!", "💰 Успех начинается с грамотного управления финансами!",
		"💰 Настройся на экономию и достигай целей!", "💰 Каждый шаг к учёту — шаг к финансовой свободе!",
		"💰 Добро пожаловать в мир финансового порядка!", "💰 Везде нужен порядок. Особенно в деньгах!",
		"💰 Время финансовой грамотности — начинаем!", "💰 Давай посмотрим, сколько ты сэкономил сегодня!",
		"💰 Контроль за финансами — путь к твоим мечтам!", "💰 Будь хозяином своих денег, а не их рабом!",
		"💰 Сегодня отличный день для достижения финансовых целей!", "💰 Финансы под контролем — стресс под нулём!",
		"💰 Даже маленькие шаги ведут к большим достижениям!", "💰 Каждый записанный расход приближает тебя к успеху!",
		"💰 Учёт финансов — как фитнес, только для кошелька!", "💰 Вложи время в учёт сегодня, пожинай успех завтра!",
		"💰 Кто владеет своими расходами, владеет своим будущим!", "💰 Давай проложим путь к твоим сбережениям!",
		"💰 Лёгкий учёт сегодня — спокойный сон завтра!", "💰 Миром правят деньги. Не забывай!",
		"💰 Умный подход к деньгам — это новый уровень жизни!", "💰 Все траты видны — будущее понятно!",
		"💰 Сэкономил время, когда начал считать деньги!", "💰 Деньги любят порядок, а порядок начинаем с учёта!",
		"💰 Давай проверим, насколько эффективно ты распоряжаешься бюджетом!", "💰 Твои деньги должны работать на тебя, а не ты на них!",
	}
	return hiText[rand.Intn(len(hiText))]
}
