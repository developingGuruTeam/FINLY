package methodsForSummary

import (
	"cachManagerApp/app/internal/methodsForAnalytic/methodsForExpenses"
	"cachManagerApp/app/internal/methodsForAnalytic/methodsForIncomeAnalys"
	"cachManagerApp/database"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// все писал в тг, это отчет , который повторяет другие и бесполезен. как по мне лучше зарезервить кнопку подо что-то ИЛИ
// ЗДЕСЬ предлагаю сделать сборник топов как например, топ по каждой категории трат, приходов, какую то статистику из бж собрать
// самый прибыльный день
// самый затратный день
// можно вообще вывести сколько операций внесено в бд приходов и расходов количество, которые это пользователи ппц как любят
// а также предлагаю убрать и неделю и месяц тут в дальнейшем меню

func AnalyseByCategoriesWeek(update tgbotapi.Update) (string, error) {

	if database.DB == nil {
		return "", fmt.Errorf("база данных не инициализирована")
	}

	analyticExpenses := methodsForExpenses.ExpensesHandler{DB: database.DB}
	analyticIncomes := methodsForIncomeAnalys.AnalyticHandler{DB: database.DB}

	totalWeekExpenses, err := analyticExpenses.ExpenseWeekAnalytic(update)
	if err != nil {
		return "", fmt.Errorf("ошибка при анализе расходов за неделю: %v", err)
	}

	totalWeekIncomes, err := analyticIncomes.IncomeWeekAnalytic(update)
	if err != nil {
		return "", fmt.Errorf("ошибка при анализе доходов за неделю: %v", err)
	}

	// Определяем категории с наибольшими расходами и доходами
	var topExpenseCategory string
	var maxExpense uint64

	var topIncomeCategory string
	var maxIncome uint64

	// Суммируем расходы и находим топовую категорию
	for category, amount := range totalWeekExpenses {
		if amount > maxExpense {
			maxExpense = amount
			topExpenseCategory = category
		}
	}
	// Суммируем доходы и находим топовую категорию
	for category, amount := range totalWeekIncomes {
		if amount > maxIncome {
			maxIncome = amount
			topIncomeCategory = category
		}
	}
	// Генерация итогового текста
	report := fmt.Sprintf("%s Аналитика доходов и расходов по категориям\n\n", update.Message.Chat.LastName)

	// Расходы по категориям
	if len(totalWeekExpenses) > 0 {
		report += "💸 Вы жадно тратили по категориям:\n"
		for category := range totalWeekExpenses {
			report += fmt.Sprintf("   ▪ %s\n", category)
		}
		report += fmt.Sprintf("\n😱 Больше всего расходов в категории: %s - %d\n", topExpenseCategory, maxExpense)
	} else {
		report += "💸 Расходов за неделю не обнаружено.\n"
	}

	report += "\n"

	// Доходы по категориям
	if len(totalWeekIncomes) > 0 {
		report += "💵 Вы безжалостно зарабатывали по категориям:\n"
		for category, _ := range totalWeekIncomes {
			report += fmt.Sprintf("   ▪ %s\n", category)
		}
		report += fmt.Sprintf("\n🤑 Больше всего доходов в категории: %s - %d\n", topIncomeCategory, maxIncome)
	} else {
		report += "💵 Доходов за неделю не обнаружено.\n"
	}
	return report, nil
}

// анализ по категориям за месяц
func AnalyseByCategoriesMonth(update tgbotapi.Update) (string, error) {

	if database.DB == nil {
		return "", fmt.Errorf("база данных не инициализирована")
	}

	analyticExpenses := methodsForExpenses.ExpensesHandler{DB: database.DB}
	analyticIncomes := methodsForIncomeAnalys.AnalyticHandler{DB: database.DB}

	totalMonthExpenses, err := analyticExpenses.ExpenseMonthAnalytic(update)
	if err != nil {
		return "", fmt.Errorf("ошибка при анализе расходов за неделю: %v", err)
	}

	totalMonthIncomes, _, err := analyticIncomes.IncomeMonthAnalytic(update)
	if err != nil {
		return "", fmt.Errorf("ошибка при анализе доходов за неделю: %v", err)
	}

	// Определяем категории с наибольшими расходами и доходами
	var topExpenseCategory string
	var maxExpense uint64

	var topIncomeCategory string
	var maxIncome uint64

	// Суммируем расходы и находим топовую категорию
	for category, amount := range totalMonthExpenses {
		if amount > maxExpense {
			maxExpense = amount
			topExpenseCategory = category
		}
	}
	// Суммируем доходы и находим топовую категорию
	for category, amount := range totalMonthIncomes {
		if amount > maxIncome {
			maxIncome = amount
			topIncomeCategory = category
		}
	}
	// Генерация итогового текста
	report := "📊 Ваш анализ за месяц по категориям:\n\n"

	// Расходы по категориям
	if len(totalMonthExpenses) > 0 {
		report += "💸 Вы жадно тратили по категориям:\n"
		for category := range totalMonthExpenses {
			report += fmt.Sprintf("   ▪ %s\n", category)
		}
		report += fmt.Sprintf("\n😱 Больше всего расходов в категории: %s - %d\n", topExpenseCategory, maxExpense)
	} else {
		report += "💸 Расходов за месяц не обнаружено.\n"
	}

	report += "\n"

	// Доходы по категориям
	if len(totalMonthIncomes) > 0 {
		report += "💵 Вы безжалостно зарабатывали по категориям:\n"
		for category := range totalMonthIncomes {
			report += fmt.Sprintf("   ▪ %s\n", category)
		}
		report += fmt.Sprintf("\n🤑 Больше всего доходов в категории: %s - %d\n", topIncomeCategory, maxIncome)
	} else {
		report += "💵 Доходов нет\n"
	}
	return report, nil
}
