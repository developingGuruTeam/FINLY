package tg_bot

import (
	methods_for_expences "cachManagerApp/app/internal/methods-for-analytic/methods-for-expenses"
	methods_for_incomes "cachManagerApp/app/internal/methods-for-analytic/methods-for-incomes"
	summary "cachManagerApp/app/internal/methods-for-analytic/methods-for-summary"
	buttons_create "cachManagerApp/app/pkg/buttons-create"
	"cachManagerApp/database"
	"fmt"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleButtonPress(bot *tgbotapi.BotAPI, update tgbotapi.Update, buttonCreator buttons_create.TelegramButtonCreator, log *slog.Logger) {
	chatID := update.Message.Chat.ID
	currency, _ := CurrencyFromChatID(chatID)

	handled := false // флажок
	switch update.Message.Text {

	// ОПИСАНИЕ КНОПОК ГЛАВНОГО МЕНЮ
	case "📥 Приход":
		incomeMenu := buttonCreator.CreateIncomeMenuButtons()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "⚙ Выберите категорию")
		msg.ReplyMarkup = incomeMenu
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send message for income:", slog.Any("error", err))
		}
		handled = true

	case "📤 Расход":
		expensesMenu := buttonCreator.CreateExpensesMenuButtons()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "⚙ Выберите категорию")
		msg.ReplyMarkup = expensesMenu
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send message for expense:", slog.Any("error", err))
		}
		handled = true

	case "🕹 Управление":
		manageMenu := buttonCreator.CreateManageMenuButtons()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "⚙ Выберите категорию")
		msg.ReplyMarkup = manageMenu
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send message for management:", slog.Any("error", err))
		}
		handled = true

	case "📊 Отчеты":
		reportMenu := buttonCreator.CreateReportsMenuButtons()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "📊 Выберите тип отчета")
		msg.ReplyMarkup = reportMenu
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send message for reports:", slog.Any("error", err))
		}
		handled = true

	case "ℹ️ Информация":
		AboutBot(bot, update.Message.Chat.ID, log)
		handled = true

	case "⚙️ Настройки":
		settingsMenu := buttonCreator.CreateSettingsMenuButtons()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "⚙ Выберите параметры")
		msg.ReplyMarkup = settingsMenu
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send message for settings:", slog.Any("error", err))
		}
		handled = true

	case "⬅ В меню":
		mainMenu := buttonCreator.CreateMainMenuButtons()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вы вернулись в главное меню")
		msg.ReplyMarkup = mainMenu
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send main menu:", slog.Any("error", err))
		}
		handled = true

	// ОПИСАНИЕ ИНЛАЙН КОММАНД

	case "/hi":
		// оставил одну инлайн команду 1 - для того что показать есть такой функционал, 2 - просто в прикол пообщаться пользователю
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, buttons_create.RandomTextForHi())
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send /help message:", slog.Any("error", err))
		}
		handled = true

	// кнопки меню НАСТРОЙКИ

	case "🎭 Изменить имя":
		clearName, _ := ClearUserNameFromChatID(chatID)
		nameText := fmt.Sprintf("Текущее имя : *%s*\n\nВведите новое имя\n_(до 32 символов)_", clearName)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, nameText)
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true) // скрываем кнопки от юзера
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send /help message:", slog.Any("error", err))
		}
		mu.Lock()
		userStates[chatID] = UserResponse{Action: "change_name"}
		mu.Unlock()
		handled = true

	case "💱 Изменить валюту":
		currencyText := fmt.Sprintf("Текущая валюта: *%s*\n\nВведите новую валюту\n_(3 символа)_\n", currency)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, currencyText)
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true) // скрываем кнопки от юзера
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send /help message:", slog.Any("error", err))
		}
		mu.Lock()
		userStates[chatID] = UserResponse{Action: "change_currency"}
		mu.Unlock()
		handled = true

	case "💫 Тарифный план":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "👷‍🔧`В разработке ...`\n")
		msg.ParseMode = "Markdown"
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send /info message:", slog.Any("error", err))
		}
		handled = true

		// кнопка меню ПРИХОД

	case "💳 Заработная плата":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите сумму заработной платы\n")
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true) // скрываем кнопки от юзера
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send /help message:", slog.Any("error", err))
		}
		mu.Lock()
		transactionStates[chatID] = TransactionResponse{Action: "salary"}
		mu.Unlock()
		handled = true

	case "🌟 Побочный доход":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите сумму дополнительного дохода\n_(подработка, фриланс)_\n")
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true) // скрываем кнопки от юзера
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send /help message:", slog.Any("error", err))
		}
		mu.Lock()
		transactionStates[chatID] = TransactionResponse{Action: "additional_income"}
		mu.Unlock()
		handled = true

	case "💼 Доход от бизнеса":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите сумму дохода от бизнеса\n")
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true) // скрываем кнопки от юзера
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send /help message:", slog.Any("error", err))
		}
		mu.Lock()
		transactionStates[chatID] = TransactionResponse{Action: "business_income"}
		mu.Unlock()
		handled = true

	case "🏦 Доход от инвестиций":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите сумму дохода от инвестиций\n_(проценты по вкладам, дивиденды)_\n")
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true) // скрываем кнопки от юзера
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send /help message:", slog.Any("error", err))
		}
		mu.Lock()
		transactionStates[chatID] = TransactionResponse{Action: "investment_income"}
		mu.Unlock()
		handled = true

	case "👮‍ Гос. выплаты":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите сумму государственных выплат\n_(пенсии, пособия)_\n")
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true) // скрываем кнопки от юзера
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send /help message:", slog.Any("error", err))
		}
		mu.Lock()
		transactionStates[chatID] = TransactionResponse{Action: "state_payments"}
		mu.Unlock()
		handled = true

	case "🏠 Продажа имущества":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите сумму от реализации имущества\n")
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true) // скрываем кнопки от юзера
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send /help message:", slog.Any("error", err))
		}
		mu.Lock()
		transactionStates[chatID] = TransactionResponse{Action: "property_sales"}
		mu.Unlock()
		handled = true

	case "⚪️ Прочие доходы":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите сумму прочих поступлений\n")
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true) // скрываем кнопки от юзера
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send /help message:", slog.Any("error", err))
		}
		mu.Lock()
		transactionStates[chatID] = TransactionResponse{Action: "other_income"}
		mu.Unlock()
		handled = true

		// кнопка меню РАСХОД

	case "🛍 Бытовые траты":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите сумму бытовых расходов\n_(еда, напитки, проезд)_\n")
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true) // скрываем кнопки от юзера
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send /help message:", slog.Any("error", err))
		}
		mu.Lock()
		transactionStates[chatID] = TransactionResponse{Action: "basic_expense"}
		mu.Unlock()
		handled = true

	case "♻️ Регулярные платежи":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите сумму регулярного платежа\n_(кредиты, налоги, аренда,\nкоммунальные платежи)_\n")
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true) // скрываем кнопки от юзера
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send /help message:", slog.Any("error", err))
		}
		mu.Lock()
		transactionStates[chatID] = TransactionResponse{Action: "regular_expense"}
		mu.Unlock()
		handled = true

	case "👘 Одежда":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите сумму трат на обновление гардероба\n")
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true) // скрываем кнопки от юзера
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send /help message:", slog.Any("error", err))
		}
		mu.Lock()
		transactionStates[chatID] = TransactionResponse{Action: "clothes"}
		mu.Unlock()
		handled = true

	case "💪 Здоровье":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите расходы на поддержание здоровья\n_(аптеки, обследования, визиты к врачам)_\n")
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true) // скрываем кнопки от юзера
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send /help message:", slog.Any("error", err))
		}
		mu.Lock()
		transactionStates[chatID] = TransactionResponse{Action: "health"}
		mu.Unlock()
		handled = true

	case "👨‍🏫 Досуг и образование":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите сумму расхода\n_(книги, подписки, курсы, хобби,\n музеи, кино, рестораны)_\n")
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true) // скрываем кнопки от юзера
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send /help message:", slog.Any("error", err))
		}
		mu.Lock()
		transactionStates[chatID] = TransactionResponse{Action: "leisure_education"}
		mu.Unlock()
		handled = true

	case "🏦 Инвестиции":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите сумму затраченную на инвестиции\n_(вклады, акции, покупка автомобилей,\nнедвижимости, предметов роcкоши)_\n")
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true) // скрываем кнопки от юзера
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send /help message:", slog.Any("error", err))
		}
		mu.Lock()
		transactionStates[chatID] = TransactionResponse{Action: "investment_expense"}
		mu.Unlock()
		handled = true

	case "⚪️ Прочие расходы":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите сумму прочих расходов\n")
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true) // скрываем кнопки от юзера
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send /help message:", slog.Any("error", err))
		}
		mu.Lock()
		transactionStates[chatID] = TransactionResponse{Action: "other_expense"}
		mu.Unlock()
		handled = true

	// кнопка меню ОТЧЕТЫ (доходы)

	case "💵 Отчет по доходам":
		incomes := buttonCreator.CreateIncomeAnalyticButtons()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите период отчета")
		msg.ReplyMarkup = incomes
		if _, err := bot.Send(msg); err != nil {
			log.Error("Failed to send main menu:", slog.Any("error", err))
		}
		handled = true

	case "📈 Отчет за день":
		analyticHandler := methods_for_incomes.AnalyticHandler{DB: database.DB} // Подключение к базе

		// Получаем данные за день
		transactions, err := analyticHandler.IncomeDayAnalytic(update)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Не удалось получить данные. Попробуйте позже.")
			_, _ = bot.Send(msg)
			log.Error("Ошибка получения данных за день:", slog.Any("error", err))
			return
		}

		// Формируем текст отчёта
		report := methods_for_incomes.GenerateDailyIncomeReport(transactions, currency)
		msg := tgbotapi.NewMessage(chatID, report)
		msg.ParseMode = "Markdown"
		_, _ = bot.Send(msg)
		handled = true

	case "📈 Отчет за неделю":
		dbConn := methods_for_incomes.AnalyticHandler{DB: database.DB}

		// Получаем данные за неделю
		incomeSummary, err := dbConn.IncomeWeekAnalytic(update)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Не удалось получить данные. Попробуйте позже.")
			_, _ = bot.Send(msg)
			log.Error("Ошибка получения данных по доходам за неделю:", slog.Any("error", err))
			return
		}

		// Генерируем текстовый отчет
		report := methods_for_incomes.GenerateWeeklyIncomeReport(incomeSummary, currency)

		// Генерируем диаграмму
		chartURL, err := methods_for_incomes.GenerateWeeklyIncomePieChartURL(incomeSummary)
		if err != nil {
			log.Error("Ошибка генерации диаграммы:", slog.Any("error", err))
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s\n(Диаграмму построить не удалось)", report))
			_, _ = bot.Send(msg)
			handled = true
			return
		}

		// Отправляем диаграмму с подписью
		imageMsg := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(chartURL))
		imageMsg.Caption = report       // Устанавливаем подпись из текста отчёта
		imageMsg.ParseMode = "Markdown" // Форматирование текста в подписи
		_, err = bot.Send(imageMsg)
		if err != nil {
			log.Info("Ошибка отправки изображения с подписью:", slog.Any("error", err))
			return
		}

		handled = true

	case "📈 Отчет за месяц":
		analyticHandler := methods_for_incomes.AnalyticHandler{DB: database.DB} // Подключение к базе

		// Получаем данные за месяц
		transactions, totalIncome, err := analyticHandler.IncomeMonthAnalytic(update)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Не удалось получить данные. Попробуйте позже.")
			_, _ = bot.Send(msg)
			log.Info("Ошибка получения данных за месяц", slog.Any("error", err))
			return
		}

		// Генерируем текстовый отчёт
		report := methods_for_incomes.GenerateMonthlyIncomeReport(transactions, currency)

		// Генерируем URL диаграммы
		chartURL, err := methods_for_incomes.GenerateIncomePieChartURL(transactions, totalIncome)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Ошибка генерации диаграммы. Попробуйте позже.")
			_, _ = bot.Send(msg)
			log.Info("Ошибка генерации графика", slog.Any("error", err))
			return
		}

		// Отправляем диаграмму с подписью
		imageMsg := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(chartURL))
		imageMsg.Caption = report       // Устанавливаем подпись из текста отчёта
		imageMsg.ParseMode = "Markdown" // Форматирование текста в подписи
		_, err = bot.Send(imageMsg)
		if err != nil {
			log.Info("Ошибка отправки изображения с подписью", slog.Any("error", err))
			return
		}

		handled = true

	// кнопка меню ОТЧЕТЫ (расходы)

	case "💸 Отчет по расходам":
		incomes := buttonCreator.CreateExpensesAnalyticButtons()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите период отчета")
		msg.ReplyMarkup = incomes
		if _, err := bot.Send(msg); err != nil {
			log.Info("Failed to send main menu", slog.Any("error", err))
		}
		handled = true

	case "📉 Отчет за день":
		dbConn := methods_for_expences.ExpensesHandler{DB: database.DB}
		expenses, err := dbConn.ExpenseDayAnalytic(update)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Не удалось получить данные. Попробуйте позже.")
			_, _ = bot.Send(msg)
			log.Info("Ошибка получения данных за день", slog.Any("error", err))
			return
		}
		report := methods_for_expences.GenerateDailyExpenseReport(expenses, currency)
		msg := tgbotapi.NewMessage(chatID, report)
		msg.ParseMode = "Markdown"
		_, _ = bot.Send(msg)
		handled = true

	case "📉 Отчет за неделю":
		dbConn := methods_for_expences.ExpensesHandler{DB: database.DB}
		expenses, err := dbConn.ExpenseWeekAnalytic(update) // Получаем данные за неделю
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Не удалось получить данные. Попробуйте позже.")
			_, _ = bot.Send(msg)
			log.Info("Ошибка получения данных по расходам за неделю", slog.Any("error", err))
			return
		}

		report := methods_for_expences.GenerateWeeklyExpensesReport(expenses, currency) // отчет

		// строим диаграмму
		chartURL, err := methods_for_expences.GenerateWeeklyExpensePieChartURL(expenses)
		if err != nil {
			log.Info("Ошибка генерации диаграммы", slog.Any("error", err))
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s\n(Диаграмму построить не удалось)", report))
			_, _ = bot.Send(msg)
			handled = true
			return
		}

		// Отправляем диаграмму с подписью
		imageMsg := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(chartURL))
		imageMsg.Caption = report       // Устанавливаем подпись из текста отчёта
		imageMsg.ParseMode = "Markdown" // Форматирование текста в подписи
		_, err = bot.Send(imageMsg)
		if err != nil {
			log.Info("Ошибка отправки изображения с подписью", slog.Any("error", err))
			return
		}

		handled = true

	case "📉 Отчет за месяц":
		dbConn := methods_for_expences.ExpensesHandler{DB: database.DB}
		expenses, err := dbConn.ExpenseMonthAnalytic(update)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Не удалось получить данные. Попробуйте позже.")
			_, _ = bot.Send(msg)
			log.Info("Ошибка получения данных по расходам за месяц", slog.Any("error", err))
			return
		}
		report := methods_for_expences.GenerateMonthlyExpensesReport(expenses, currency)

		// строим диаграмму
		chartURL, err := methods_for_expences.GenerateExpensePieChartURL(expenses)
		if err != nil {
			log.Info("Ошибка генерации диаграммы", slog.Any("error", err))
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s\n(Диаграмму построить не удалось)", report))
			_, _ = bot.Send(msg)
			handled = true
			return
		}

		// Отправляем диаграмму с подписью
		imageMsg := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(chartURL))
		imageMsg.Caption = report       // Устанавливаем подпись из текста отчёта
		imageMsg.ParseMode = "Markdown" // Форматирование текста в подписи
		_, err = bot.Send(imageMsg)
		if err != nil {
			log.Info("Ошибка отправки изображения с подписью", slog.Any("error", err))
			return
		}

		handled = true

	// кнопка меню УПРАВЛЕНИЕ

	case "🛎 Напоминание":
		command := "🛎 Напоминание"
		PushOnAnalyticButton(bot, update, buttonCreator, command, log)
		handled = true

	case "🗓 Подписки":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "👷‍🔧`В разработке ...`\n")
		msg.ParseMode = "Markdown"
		if _, err := bot.Send(msg); err != nil {
			log.Info("Failed to send /info message: ", slog.Any("error", err))
		}
		handled = true

	// кнопки меню внутри Отчетов
	case "🧑‍💻 Аналитика":
		analyse := buttonCreator.CreateSuperAnalyticButtons()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите категорию аналитики")
		msg.ReplyMarkup = analyse
		if _, err := bot.Send(msg); err != nil {
			log.Info("Failed to send main menu:", slog.Any("error", err))
		}
		handled = true

	case "💲Анализ за неделю":
		command := "💲Анализ за неделю"
		PushOnAnalyticButton(bot, update, buttonCreator, command, log)
		handled = true

	case "💰Анализ за месяц":
		command := "💰Анализ за месяц"
		PushOnAnalyticButton(bot, update, buttonCreator, command, log)
		handled = true

	case "🧮 Статистика":
		dbConn := database.DB
		userID := update.Message.From.ID
		report := summary.GenerateStatisticsReport(userID, dbConn)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, report)
		msg.ParseMode = "Markdown"
		if _, err := bot.Send(msg); err != nil {
			log.Info("Ошибка отправки сообщения:", slog.Any("error", err))
		}

		handled = true

	case "⚖️ Cальдо":
		command := "⚖️ Cальдо"
		PushOnAnalyticButton(bot, update, buttonCreator, command, log)
		handled = true

	case "💲Сальдо за неделю":
		command := "💲Сальдо за неделю"
		PushOnAnalyticButton(bot, update, buttonCreator, command, log)
		handled = true

	case "💰Сальдо за месяц":
		command := "💰Сальдо за месяц"
		PushOnAnalyticButton(bot, update, buttonCreator, command, log)
		handled = true

	case "👨‍🔬 Экспертная аналитика":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "👷‍🔧`В разработке ...`\n")
		msg.ParseMode = "Markdown"
		if _, err := bot.Send(msg); err != nil {
			log.Info("Failed to send main menu:", slog.Any("error", err))
		}
		handled = true

	// кнопки меню внутри Управления - Напоминания

	case "🔁 Регулярный платёж":
		command := "🔁 Регулярный платёж"
		PushOnAnalyticButton(bot, update, buttonCreator, command, log)
		handled = true

	case "🎯 Накопления":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "👷‍🔧`В разработке ...`\n")
		msg.ParseMode = "Markdown"
		if _, err := bot.Send(msg); err != nil {
			log.Info("Failed to send /info message:", slog.Any("error", err))
		}
		handled = true

	// предлагаю сделать напоминание настраиваемое прям когда человек хочет) одноразовое хоть через 3 дня хоть через 333 дня
	case "🔂 Разовый платеж":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "👷‍🔧`В разработке ...`\n")
		msg.ParseMode = "Markdown"
		if _, err := bot.Send(msg); err != nil {
			log.Info("Failed to send /info message:", slog.Any("error", err))
		}
		handled = true
	}

	// Если команда или кнопка не обработаны, отправляем сообщение об ошибке
	if !handled {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🚫 Неизвестная команда. Повторите запрос.")
		if _, err := bot.Send(msg); err != nil {
			log.Info("Failed to send unknown command message:", slog.Any("error", err))
		}
	}
}
