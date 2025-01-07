package methods_for_summary

import (
	"cachManagerApp/app/db/models"
	expenses "cachManagerApp/app/internal/methods-for-analytic/methods-for-expenses"
	incomes "cachManagerApp/app/internal/methods-for-analytic/methods-for-incomes"
	"cachManagerApp/database"
	"fmt"
	"math/rand"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	summary = models.Summary{}
)

// анализ сальдо за неделю
func AnalyseBySaldoWeek(update tgbotapi.Update) (models.Summary, error) {
	analyticExpenses := expenses.ExpensesHandler{DB: database.DB}
	analyticIncomes := incomes.AnalyticHandler{DB: database.DB}
	var summary models.Summary

	if database.DB == nil {
		return models.Summary{}, fmt.Errorf("база данных не инициализирована в сальдо за неделю")
	}

	totalWeekExpenses, err := analyticExpenses.ExpenseWeekAnalytic(update)
	if err != nil {
		return models.Summary{}, fmt.Errorf("ошибка при анализе расходов за неделю: %v", err)
	}

	totalWeekIncomes, err := analyticIncomes.IncomeWeekAnalytic(update)
	if err != nil {
		return models.Summary{}, fmt.Errorf("ошибка при анализе доходов за неделю: %v", err)
	}
	fmt.Println(summary.TotalIncome)
	// подсчет общих доходов
	for category, amount := range totalWeekIncomes {
		summary.TotalIncome += amount
		summary.IncomeCategories = append(summary.IncomeCategories, models.CategorySummary{Category: category, Amount: amount})

		if amount > summary.TopIncome.Amount {
			summary.TopIncome = models.CategorySummary{
				Category: category,
				Amount:   amount,
			}
		}
	}
	fmt.Print(summary.TotalIncome)
	for category, amount := range totalWeekExpenses {
		summary.TotalExpense += amount
		summary.ExpenseCategories = append(summary.ExpenseCategories, models.CategorySummary{Category: category, Amount: amount})

		if amount > summary.TopExpense.Amount {
			summary.TopExpense = models.CategorySummary{
				Category: category,
				Amount:   amount,
			}
		}
	}

	summary.Profit = int64(summary.TotalIncome) - int64(summary.TotalExpense)
	return summary, nil
}

func GenerateWeeklySaldoReport(sum models.Summary, currency string) string {
	report := "📊 *Сальдо за неделю*\n"

	// Итоговая прибыль или убыток
	if sum.Profit >= 0 {
		report += fmt.Sprintf("\n🟢 Баланс положительный *%d* %s\n", sum.Profit, currency)
	} else {
		report += fmt.Sprintf("\n🔴 Баланс отрицательный *%d* %s\n", -sum.Profit, currency)
	}

	// топ расход
	if len(sum.ExpenseCategories) > 0 {
		report += fmt.Sprintf("\n💸 Наибольшие расходы в категории:\n%s %d %s \n", sum.TopExpense.Category, sum.TopExpense.Amount, currency)
	} else {
		report += "   ▪ Расходов нет.\n"
	}

	// топ доход
	if len(sum.IncomeCategories) > 0 {
		report += fmt.Sprintf("\n💵 Наибольшие доходы в категории:\n%s %d %s\n", sum.TopIncome.Category, sum.TopIncome.Amount, currency)
	} else {
		report += "   ▪ Доходов нет.\n"
	}

	if sum.Profit >= 0 {
		report += fmt.Sprintf("\n\n💡`%s`\n ", randomAdvicePositive())
	} else {
		report += fmt.Sprintf("\n\n💡`%s`\n", randomAdviceNegative())
	}
	return report
}

// анализ сальдо за месяц
func AnalyseBySaldoMonth(update tgbotapi.Update) (models.Summary, error) {
	analyticExpenses := expenses.ExpensesHandler{DB: database.DB}
	analyticIncomes := incomes.AnalyticHandler{DB: database.DB}
	var summary models.Summary
	if database.DB == nil {
		return models.Summary{}, fmt.Errorf("ошибка подключения к БД в аналитике сальдо")
	}

	totalExpenses, err := analyticExpenses.ExpenseMonthAnalytic(update)
	if err != nil {
		return models.Summary{}, fmt.Errorf("ошибка в вычислении расходов")
	}
	totalIncomes, _, err := analyticIncomes.IncomeMonthAnalytic(update)
	if err != nil {
		return models.Summary{}, fmt.Errorf("ошибка в вычислении доходов")
	}

	for category, amount := range totalExpenses {
		summary.TotalExpense += amount
		summary.ExpenseCategories = append(summary.ExpenseCategories, models.CategorySummary{Category: category, Amount: amount})
		if amount > summary.TopExpense.Amount {
			summary.TopExpense = models.CategorySummary{
				Category: category,
				Amount:   amount}
		}
	}

	for category, amount := range totalIncomes {
		summary.TotalIncome += amount
		summary.IncomeCategories = append(summary.IncomeCategories, models.CategorySummary{Category: category, Amount: amount})
		if amount > summary.TopIncome.Amount {
			summary.TopIncome = models.CategorySummary{Category: category, Amount: amount}
		}
	}

	summary.Profit = int64(summary.TotalIncome) - int64(summary.TotalExpense)
	return summary, nil
}

func GenerateMonthlySaldoReport(sum models.Summary, currency string) string {
	report := "📊 *Сальдо за месяц*\n"

	if sum.Profit >= 0 {
		report += fmt.Sprintf("\n🟢 Баланс положительный *%d* %s\n ", sum.Profit, currency)
	} else {
		report += fmt.Sprintf("\n🔴 Баланс отрицательный *%d* %s\n", -sum.Profit, currency)
	}

	// топ расход
	if len(sum.ExpenseCategories) > 0 {
		report += fmt.Sprintf("\n💸 Наибольшие расходы в категории:\n%s %d %s\n", sum.TopExpense.Category, sum.TopExpense.Amount, currency)
	} else {
		report += "   ▪ Расходов нет.\n"
	}

	// топ доход
	if len(sum.IncomeCategories) > 0 {
		report += fmt.Sprintf("\n💵 Наибольшие доходы в категории:\n%s %d %s \n", sum.TopIncome.Category, sum.TopIncome.Amount, currency)
	} else {
		report += "   ▪ Доходов нет.\n"
	}

	if sum.Profit >= 0 {
		report += fmt.Sprintf("\n\n💡`%s`\n ", randomAdvicePositive())
	} else {
		report += fmt.Sprintf("\n\n💡`%s`\n", randomAdviceNegative())
	}
	return report
}

func randomAdviceNegative() string {
	txt := [...]string{
		"Иногда стоит задумываться о своих тратах. Обратите внимание на свои финансовые показатели.",
		"Возможно, у Вас были веские основания для таких затрат. Однако, нужно быть более экономным.",
		"Проверьте, куда уходят ваши деньги. Возможно, некоторые траты можно сократить.",
		"Попробуйте планировать бюджет на неделю вперед, чтобы избежать лишних расходов.",
		"Иногда полезно записывать все траты для лучшего контроля финансов.",
		"Слишком много лишних расходов? Пересмотрите приоритеты.",
		"Помните, что даже маленькие сбережения могут привести к большим результатам.",
		"Может, стоит задуматься о копилке для небольших, но важных целей?",
		"Избегайте импульсивных покупок — это главный враг финансовой стабильности.",
		"Если расходы превысили доходы, попробуйте временно урезать ненужные траты.",
		"Запомните: разумный контроль бюджета — ключ к успеху.",
		"Оцените, действительно ли каждая покупка была необходимой.",
		"Сделайте небольшой анализ за месяц: где можно было бы сэкономить?",
		"Не забывайте оставлять немного денег для чрезвычайных ситуаций.",
		"Попробуйте использовать скидочные карты или акции для экономии.",
		"Возможно, стоит отказаться от необязательных подписок и услуг.",
		"Попробуйте задать себе вопрос: это желание или необходимость?",
		"Контроль за мелкими расходами — первый шаг к большому результату.",
		"Проанализируйте свои траты. Возможно, можно отказаться от части привычек.",
		"Увеличение расходов — это нормально, но не забывайте о своих финансовых целях.",
	}
	return txt[rand.Intn(len(txt))]
}

func randomAdvicePositive() string {
	txt := [...]string{
		"Вы прекрасно ведете свой финансовый учет!",
		"Так держать!",
		"Ваши финансовые успехи заслуживают похвалы!",
		"Вы демонстрируете отличные навыки планирования бюджета.",
		"Ваши усилия окупаются — продолжайте в том же духе!",
		"Вы уверенно двигаетесь к своим финансовым целям.",
		"Ваши траты под контролем — это заслуживает уважения.",
		"Экономия без стресса — это ваш подход, и он работает!",
		"Вы находите баланс между расходами и доходами — это успех!",
		"Ваш бюджет — пример для подражания.",
		"Вы делаете всё правильно, чтобы достичь финансовой стабильности.",
		"Ваш подход к экономии впечатляет.",
		"Молодец! Вы экономите деньги и работаете над своими целями.",
		"Вы демонстрируете отличную финансовую дисциплину.",
		"Ваш контроль над бюджетом — это ключ к свободе.",
		"Каждая ваша экономия приближает вас к финансовой независимости.",
		"Вы сделали огромный шаг к своим финансовым мечтам.",
		"Вам удается находить разумный баланс между тратами и сбережениями.",
		"Продолжайте в том же духе, и финансовый успех не заставит себя ждать.",
		"Ваши усилия по контролю бюджета вдохновляют!",
	}
	return txt[rand.Intn(len(txt))]
}
