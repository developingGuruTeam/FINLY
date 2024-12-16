package methodsForSummary

import (
	"cachManagerApp/app/db/models"
	"cachManagerApp/app/internal/methodsForAnalytic/methodsForExpenses"
	"cachManagerApp/app/internal/methodsForAnalytic/methodsForIncomeAnalys"
	"cachManagerApp/database"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	summary = models.Summary{}
)

// анализ сальдо за неделю
func AnalyseBySaldoWeek(update tgbotapi.Update) (models.Summary, error) {
	analyticExpenses := methodsForExpenses.ExpensesHandler{DB: database.DB}
	analyticIncomes := methodsForIncomeAnalys.AnalyticHandler{DB: database.DB}

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

func GenerateWeeklySaldoReport(sum models.Summary) string {
	report := "📊 Ваш анализ за неделю:\n\n"

	// расходы
	report += "💸 Расходы по категориям:\n"
	if len(sum.ExpenseCategories) > 0 {
		for _, category := range sum.ExpenseCategories {
			report += fmt.Sprintf("   ▪ %s: %d\n", category.Category, category.Amount)
		}
		report += fmt.Sprintf("\n🔴 Больше всего расходов в категории: %s (%d)\n", sum.TopExpense.Category, sum.TopExpense.Amount)
	} else {
		report += "   ▪ Нет расходов за неделю.\n"
	}

	// Доходы
	report += "\n💵 Доходы по категориям:\n"
	if len(sum.IncomeCategories) > 0 {
		for _, category := range sum.IncomeCategories {
			report += fmt.Sprintf("   ▪ %s: %d\n", category.Category, category.Amount)
		}
		report += fmt.Sprintf("\n🟢 Больше всего доходов в категории: %s (%d)\n", sum.TopIncome.Category, sum.TopIncome.Amount)
	} else {
		report += "   ▪ Нет доходов за неделю.\n"
	}

	// Итоговая прибыль или убыток
	if sum.Profit >= 0 {
		report += fmt.Sprintf("\n✅ Итоговая прибыль за неделю: %d\n", sum.Profit)
	} else {
		report += fmt.Sprintf("\n❌ Итоговый убыток за неделю: %d\n", -sum.Profit)
	}

	return report
}

// анализ сальдо за месяц
func AnalyseBySaldoMonth(update tgbotapi.Update) (models.Summary, error) {
	analyticExpenses := methodsForExpenses.ExpensesHandler{DB: database.DB}
	analyticIncomes := methodsForIncomeAnalys.AnalyticHandler{DB: database.DB}

	if database.DB == nil {
		return models.Summary{}, fmt.Errorf("ошибка подключения к БД в аналитике сальдо")
	}

	totalExpenses, err := analyticExpenses.ExpenseMonthAnalytic(update)
	if err != nil {
		return models.Summary{}, fmt.Errorf("ошибка в вычислении расходов")
	}
	totalIncomes, err := analyticIncomes.IncomeMonthAnalytic(update)
	if err != nil {
		return models.Summary{}, fmt.Errorf("ошибка в вычислении доходов")
	}

	for category, amount := range totalExpenses {
		summary.TotalExpense += amount
		summary.ExpenseCategories = append(summary.ExpenseCategories, models.CategorySummary{category, amount})
		if amount > summary.TopExpense.Amount {
			summary.TopExpense = models.CategorySummary{category, amount}
		}
	}

	for category, amount := range totalIncomes {
		summary.TotalIncome += amount
		summary.IncomeCategories = append(summary.IncomeCategories, models.CategorySummary{category, amount})
		if amount > summary.TopIncome.Amount {
			summary.TopIncome = models.CategorySummary{category, amount}
		}
	}

	summary.Profit = int64(summary.TotalIncome) - int64(summary.TotalExpense)
	return summary, nil
}

func GenerateMonthlySaldoReport(sum models.Summary) string {
	report := "📊 Ваш анализ за месяц:\n\n"

	// расходы
	report += "💸 Расходы по категориям:\n"
	if len(sum.ExpenseCategories) > 0 {
		for _, category := range sum.ExpenseCategories {
			report += fmt.Sprintf("   ▪ %s: %d\n", category.Category, category.Amount)
		}
		report += fmt.Sprintf("\n🔴 Больше всего расходов в категории: %s (%d)\n", sum.TopExpense.Category, sum.TopExpense.Amount)
	} else {
		report += "   ▪ Нет расходов за месяц.\n"
	}

	// Доходы
	report += "\n💵 Доходы по категориям:\n"
	if len(sum.IncomeCategories) > 0 {
		for _, category := range sum.IncomeCategories {
			report += fmt.Sprintf("   ▪ %s: %d\n", category.Category, category.Amount)
		}
		report += fmt.Sprintf("\n🟢 Больше всего доходов в категории: %s (%d)\n", sum.TopIncome.Category, sum.TopIncome.Amount)
	} else {
		report += "   ▪ Нет доходов за месяц.\n"
	}

	// Итоговая прибыль или убыток
	if sum.Profit >= 0 {
		report += fmt.Sprintf("\n✅ Итоговая прибыль за месяц: %d\n", sum.Profit)
	} else {
		report += fmt.Sprintf("\n❌ Итоговый убыток за месяц: %d\n", -sum.Profit)
	}

	return report
}
