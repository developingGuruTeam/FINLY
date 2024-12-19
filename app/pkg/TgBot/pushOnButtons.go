package TgBot

import (
	"cachManagerApp/app/internal/methodsForAnalytic/methodsForExpenses"
	"cachManagerApp/app/internal/methodsForAnalytic/methodsForIncomeAnalys"
	"cachManagerApp/app/pkg/logger"
	"cachManagerApp/database"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sync"
)

type TransactionResponse struct {
	Action string `json:"action"`
}

type UserResponse struct {
	Action string `json:"action"`
}

var (
	log               = logger.GetLogger()
	userStates        = make(map[int64]UserResponse)        // мапа для хранения состояния пользователей
	mu                sync.Mutex                            // мьютекс для синхронизации доступа к мапе
	transactionStates = make(map[int64]TransactionResponse) // мапа для хранения состояния транзакций

)

// обработка нажатий на кнопки (команда приходит сюда)
func PushOnButton(bot *tgbotapi.BotAPI, update tgbotapi.Update, buttonCreator TelegramButtonCreator) {
	if update.Message != nil {
		// чат ID наполняется
		chatID := update.Message.Chat.ID
		mu.Lock()
		val2, ok2 := transactionStates[chatID]
		val, ok := userStates[chatID]
		mu.Unlock()

		if ok2 && val2.Action != "" {
			handleTransactionAction(bot, update, val2)
			return
		}

		// если в ней лежит ключ, то переходит к действию, если нет, то ждет отклика
		if ok && val.Action != "" {
			handleUserAction(bot, update, val)
			return
		}
		handleButtonPress(bot, update, buttonCreator)
	}
}

func handleButtonPress(bot *tgbotapi.BotAPI, update tgbotapi.Update, buttonCreator TelegramButtonCreator) {
	chatID := update.Message.Chat.ID
	currency, _ := CurrencyFromChatID(chatID)

	handled := false
	switch update.Message.Text {

	// ОПИСАНИЕ КНОПОК МЕНЮ
	case "📥 Приход":
		incomeMenu := buttonCreator.CreateIncomeMenuButtons()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "⚙ Выберите категорию")
		msg.ReplyMarkup = incomeMenu
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send message for income: %v", err)
		}
		handled = true

	case "📤 Расход":
		expensesMenu := buttonCreator.CreateExpensesMenuButtons()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "⚙ Выберите категорию")
		msg.ReplyMarkup = expensesMenu
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send message for expense: %v", err)
		}
		handled = true

	case "📊 Отчеты":
		reportMenu := buttonCreator.CreateReportsMenuButtons()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "📊 Выберите тип отчета")
		msg.ReplyMarkup = reportMenu
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send message for reports: %v", err)
		}
		handled = true

	case "⚙️ Настройки":
		settingsMenu := buttonCreator.CreateSettingsMenuButtons()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "⚙ Выберите параметры")
		msg.ReplyMarkup = settingsMenu
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send message for settings: %v", err)
		}
		handled = true

	case "⬅ В меню":
		mainMenu := buttonCreator.CreateMainMenuButtons()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вы вернулись в главное меню")
		msg.ReplyMarkup = mainMenu
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send main menu: %v", err)
		}
		handled = true

	// ОПИСАНИЕ ИНЛАЙН КОММАНД
	case "/info":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "📍 Бот предназначен для:\n ▪ Ведения учета доходов и расходов\n ▪ Создания отчетов по различным критериям\n ▪ Экономического анализа")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send /info message: %v", err)
		}
		handled = true
	// дописать нормальный хэлп!!!!!!
	case "/help":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "📌 Команды бота:\n/info - Информация о боте\n/help - Помощь по использованию бота")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send /help message: %v", err)
		}
		handled = true

	case "/hi":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, randomTextForHi()) // дописать нормальный хэлп
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send /help message: %v", err)
		}
		handled = true

	// кнопки меню НАСТРОЙКИ
	case "🎭 Изменить имя":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите Ваше новое имя")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send /help message: %v", err)
		}
		mu.Lock()
		userStates[chatID] = UserResponse{Action: "change_name"}
		mu.Unlock()
		handled = true

	case "💱 Изменить валюту":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите валюту")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send /help message: %v", err)
		}
		mu.Lock()
		userStates[chatID] = UserResponse{Action: "change_currency"}
		mu.Unlock()
		handled = true

	case "💫 Тарифный план":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "👷‍🔧`В разработке ...`\n\n`Ожидаемая дата выхода обновления 20.01.2025`")
		msg.ParseMode = "Markdown"
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send /info message: %v", err)
		}
		handled = true

		// приходы
	case "💳 Заработная плата":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите сумму заработной платы.\n\nЧерез запятую можно добавить комментарий")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send /help message: %v", err)
		}
		mu.Lock()
		transactionStates[chatID] = TransactionResponse{Action: "salary"}
		mu.Unlock()
		handled = true

	case "💱 Побочный доход":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите сумму дополнительного дохода\n(подработка, фриланс).\n\nЧерез запятую можно добавить комментарий")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send /help message: %v", err)
		}
		mu.Lock()
		transactionStates[chatID] = TransactionResponse{Action: "additional_income"}
		mu.Unlock()
		handled = true

	case "😎 Доход от бизнеса":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите сумму дохода от бизнеса.\n\nЧерез запятую можно добавить комментарий")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send /help message: %v", err)
		}
		mu.Lock()
		transactionStates[chatID] = TransactionResponse{Action: "business_income"}
		mu.Unlock()
		handled = true

	case "🏦 Доход от инвестиций":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите сумму дохода от инвестиций\n(проценты по вкладам, дивиденды).\n\nЧерез запятую можно добавить комментарий")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send /help message: %v", err)
		}
		mu.Lock()
		transactionStates[chatID] = TransactionResponse{Action: "investment_income"}
		mu.Unlock()
		handled = true

	case "👮‍ Гос. выплаты":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите сумму государственных выплат\n(пенсии, субсидии).\n\nЧерез запятую можно добавить комментарий")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send /help message: %v", err)
		}
		mu.Lock()
		transactionStates[chatID] = TransactionResponse{Action: "state_payments"}
		mu.Unlock()
		handled = true

	case "🏠 Продажа имущества":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите сумму продажи имущества.\n\nЧерез запятую можно добавить комментарий")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send /help message: %v", err)
		}
		mu.Lock()
		transactionStates[chatID] = TransactionResponse{Action: "property_sales"}
		mu.Unlock()
		handled = true

	case "⚪️ Прочие доходы":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите сумму прочих поступлений.\n\nЧерез запятую можно добавить комментарий")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send /help message: %v", err)
		}
		mu.Lock()
		transactionStates[chatID] = TransactionResponse{Action: "other_income"}
		mu.Unlock()
		handled = true

		// расходные операции
	case "🛍 Бытовые траты":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите сумму базовых трат\n(еда, напитки, проезд).\n\nЧерез запятую можно добавить комментарий")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send /help message: %v", err)
		}
		mu.Lock()
		transactionStates[chatID] = TransactionResponse{Action: "basic_expense"}
		mu.Unlock()
		handled = true

	case "🫡 Регулярные платежи":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите сумму регулярного платежа\n(кредиты, налоги, аренда,\nкоммунальные платежи).\n\nЧерез запятую можно добавить комментарий")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send /help message: %v", err)
		}
		mu.Lock()
		transactionStates[chatID] = TransactionResponse{Action: "regular_expense"}
		mu.Unlock()
		handled = true

	case "👘 Одежда":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите сумму трат на обновление гардероба.\n\nЧерез запятую можно добавить комментарий")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send /help message: %v", err)
		}
		mu.Lock()
		transactionStates[chatID] = TransactionResponse{Action: "clothes"}
		mu.Unlock()
		handled = true

	case "💪 Здоровье":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите расходы на поддержание здоровья\n(аптеки, обследования, визиты к врачам).\n\nЧерез запятую можно добавить комментарий")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send /help message: %v", err)
		}
		mu.Lock()
		transactionStates[chatID] = TransactionResponse{Action: "health"}
		mu.Unlock()
		handled = true

	case "👨‍🏫 Досуг и образование":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите сумму расхода\n(книги, подписки, курсы, хобби,\n музеи, кино, рестораны).\n\nЧерез запятую можно добавить комментарий")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send /help message: %v", err)
		}
		mu.Lock()
		transactionStates[chatID] = TransactionResponse{Action: "leisure_education"}
		mu.Unlock()
		handled = true

	case "🚀 Инвестиции":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите сумму затраченную на инвестиции\n(вклады, акции, автомобили,\nнедвижимость, предметы роcкоши).\n\nЧерез запятую можно добавить комментарий")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send /help message: %v", err)
		}
		mu.Lock()
		transactionStates[chatID] = TransactionResponse{Action: "investment_expense"}
		mu.Unlock()
		handled = true

	case "⚪️ Прочие расходы":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите сумму прочих расходов\n\nЧерез запятую можно добавить комментарий")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send /help message: %v", err)
		}
		mu.Lock()
		transactionStates[chatID] = TransactionResponse{Action: "other_expense"}
		mu.Unlock()
		handled = true

	// ОТЧЕТ ДОХОДЫ
	case "💵 Отчет по доходам":
		incomes := buttonCreator.CreateIncomeAnalyticButtons()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите период отчета")
		msg.ReplyMarkup = incomes
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send main menu: %v", err)
		}
		handled = true

	case "📈 Отчет за день":
		analyticHandler := methodsForIncomeAnalys.AnalyticHandler{DB: database.DB} // Подключение к базе

		// Получаем данные за день
		transactions, err := analyticHandler.IncomeDayAnalytic(update)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Не удалось получить данные. Попробуйте позже.")
			_, _ = bot.Send(msg)
			log.Printf("Ошибка получения данных за день: %v", err)
			return
		}

		// Формируем текст отчёта
		report := methodsForIncomeAnalys.GenerateDailyIncomeReport(transactions, currency)
		msg := tgbotapi.NewMessage(chatID, report)
		_, _ = bot.Send(msg)
		handled = true

	case "📈 Отчет за неделю":
		analyticHandler := methodsForIncomeAnalys.AnalyticHandler{DB: database.DB} // Подключение к базе

		// Получаем данные за день
		transactions, err := analyticHandler.IncomeWeekAnalytic(update)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Не удалось получить данные. Попробуйте позже.")
			_, _ = bot.Send(msg)
			log.Printf("Ошибка получения данных за день: %v", err)
			return
		}

		// Формируем текст отчёта
		report := methodsForIncomeAnalys.GenerateWeeklyIncomeReport(transactions, currency)
		msg := tgbotapi.NewMessage(chatID, report)
		_, _ = bot.Send(msg)
		handled = true

	case "📈 Отчет за месяц":
		analyticHandler := methodsForIncomeAnalys.AnalyticHandler{DB: database.DB} // Подключение к базе

		// Получаем данные за месяц
		transactions, totalIncome, err := analyticHandler.IncomeMonthAnalytic(update)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Не удалось получить данные. Попробуйте позже.")
			_, _ = bot.Send(msg)
			log.Printf("Ошибка получения данных за месяц: %v", err)
			return
		}

		// Генерируем текстовый отчёт
		report := methodsForIncomeAnalys.GenerateMonthlyIncomeReport(transactions, currency)

		// Генерируем URL диаграммы
		chartURL, err := methodsForIncomeAnalys.GenerateIncomePieChartURL(transactions, totalIncome)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Ошибка генерации диаграммы. Попробуйте позже.")
			_, _ = bot.Send(msg)
			log.Printf("Ошибка генерации графика: %v", err)
			return
		}

		// Отправляем диаграмму
		imageMsg := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(chartURL))
		_, _ = bot.Send(imageMsg)
		// Отправляем текстовый отчёт
		textMsg := tgbotapi.NewMessage(chatID, report)
		_, _ = bot.Send(textMsg)

		handled = true

	// ОТЧЕТ ПО РАСХОДАМ
	case "💸 Отчет по расходам":
		incomes := buttonCreator.CreateExpensesAnalyticButtons()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите период отчета")
		msg.ReplyMarkup = incomes
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send main menu: %v", err)
		}
		handled = true

	case "📉 Отчет за день":
		dbConn := methodsForExpenses.ExpensesHandler{DB: database.DB}
		expenses, err := dbConn.ExpenseDayAnalytic(update)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Не удалось получить данные. Попробуйте позже.")
			_, _ = bot.Send(msg)
			log.Printf("Ошибка получения данных за день: %v", err)
			return
		}
		report := methodsForExpenses.GenerateDailyExpenseReport(expenses, currency)
		msg := tgbotapi.NewMessage(chatID, report)
		_, _ = bot.Send(msg)
		handled = true

	case "📉 Отчет за неделю":
		dbConn := methodsForExpenses.ExpensesHandler{DB: database.DB}
		expenses, err := dbConn.ExpenseWeekAnalytic(update)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Не удалось получить данные. Попробуйте позже.")
			_, _ = bot.Send(msg)
			log.Printf("Ошибка получения данных по расходам за неделю: %v", err)
			return
		}
		report := methodsForExpenses.GenerateWeeklyExpensesReport(expenses, currency)
		msg := tgbotapi.NewMessage(chatID, report)
		_, _ = bot.Send(msg)
		handled = true

	case "📉 Отчет за месяц":
		dbConn := methodsForExpenses.ExpensesHandler{DB: database.DB}
		expenses, err := dbConn.ExpenseMonthAnalytic(update)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Не удалось получить данные. Попробуйте позже.")
			_, _ = bot.Send(msg)
			log.Printf("Ошибка получения данных по расходам за месяц: %v", err)
			return
		}
		report := methodsForExpenses.GenerateMonthlyExpensesReport(expenses, currency)

		// Генерируем диаграмму
		chartURL, err := methodsForExpenses.GenerateExpensePieChartURL(expenses)
		if err != nil {
			log.Printf("Ошибка генерации диаграммы: %v", err)
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s\n(Диаграмму построить не удалось)", report))
			_, _ = bot.Send(msg)
			handled = true
			return
		}

		// Отправляем изображение диаграммы
		imageMsg := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(chartURL))
		_, _ = bot.Send(imageMsg)
		// Отправляем текстовый отчёт
		msg := tgbotapi.NewMessage(chatID, report)
		_, _ = bot.Send(msg)

		handled = true

	// аналитика
	case "🧑‍💻 Аналитика":
		analyse := buttonCreator.CreateSuperAnalyticButtons()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите категорию аналитики")
		msg.ReplyMarkup = analyse
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send main menu: %v", err)
		}
		handled = true

	case "🛍 Анализ категорий":
		command := "🛍 Анализ категорий"
		PushOnAnalyticButton(bot, update, buttonCreator, command)
		handled = true

	case "💲 Анализ за неделю":
		command := "💲 Анализ за неделю"
		PushOnAnalyticButton(bot, update, buttonCreator, command)
		handled = true

	case "💰 Анализ месяц":
		command := "💰 Анализ за месяц"
		PushOnAnalyticButton(bot, update, buttonCreator, command)
		handled = true

	case "🤑 Cальдо":
		command := "сальдо"
		PushOnAnalyticButton(bot, update, buttonCreator, command)
		handled = true

	case "💲Сальдо за неделю":
		command := "💲Сальдо за неделю"
		PushOnAnalyticButton(bot, update, buttonCreator, command)
		handled = true

	case "💰Сальдо за месяц":
		command := "💰Сальдо за месяц"
		PushOnAnalyticButton(bot, update, buttonCreator, command)
		handled = true

	case "👨‍🔬 Экспертная аналитика":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "👷‍🔧`В разработке ...`\n\n`Ожидаемая дата выхода обновления 20.01.2025`")
		msg.ParseMode = "Markdown"
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send main menu: %v", err)
		}
		handled = true
	}

	// Если команда или кнопка не обработаны, отправляем сообщение об ошибке
	if !handled {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🚫 Неизвестная команда. Повторите запрос.")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send unknown command message: %v", err)
		}
	}
}
