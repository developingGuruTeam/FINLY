package methodsForExpenses

import (
	"cachManagerApp/app/db/models"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"time"
)

type ExpensesHandler struct {
	DB *gorm.DB
}

// хендлер расходов
//go:generate mockery --name=ExpenseAnalyticHandler --output=../tests/mocks --with-expecter
type ExpenseAnalyticHandler interface {
	ExpenseDayAnalytic(update tgbotapi.Update) ([]models.Transactions, error)
	ExpenseWeekAnalytic(update tgbotapi.Update) (map[string]uint64, error)
	ExpenseMonthAnalytic(update tgbotapi.Update) (map[string]uint64, error)
}

// расход за день
func (exp *ExpensesHandler) ExpenseDayAnalytic(update tgbotapi.Update) ([]models.Transactions, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	endOfDay := startOfDay.Add(24 * time.Hour)

	var transactions []models.Transactions
	err := exp.DB.Where("telegram_id = ? AND operation_type = ? AND created_at BETWEEN ? AND ?",
		update.Message.Chat.ID, false, startOfDay, endOfDay).Find(&transactions).Error
	if err != nil {
		return nil, fmt.Errorf("error getting transactions: %v", err)
	}

	return transactions, nil
}

func GenerateDailyExpenseReport(expenses []models.Transactions) string {
	if len(expenses) == 0 {
		return "📉 Сегодня у вас не было расходов."
	}

	report := "📉 Отчёт за день:\n\n"
	var totalExpenses uint64

	for _, exp := range expenses {
		report += fmt.Sprintf("▪ Категория: %s\n", exp.Category)
		report += fmt.Sprintf("   Сумма: %d\n", exp.Quantities)
		if exp.Description != "" {
			report += fmt.Sprintf("   Комментарий: %s\n", exp.Description)
		}
		report += "\n"
		totalExpenses += exp.Quantities
	}
	report += fmt.Sprintf("💸 Итого расходов за день: %d\n", totalExpenses)
	return report
}

// расход за неделю
func (exp *ExpensesHandler) ExpenseWeekAnalytic(update tgbotapi.Update) (map[string]uint64, error) {
	now := time.Now()
	startDay := now.AddDate(0, 0, -7)
	endDay := now

	var result []struct {
		Category string
		Value    uint64
	}

	err := exp.DB.Model(&models.Transactions{}).
		Select("category, SUM (quantities) as value").
		Where("telegram_id = ? AND operation_type = ? AND created_at BETWEEN ? AND ?",
			update.Message.Chat.ID, false, startDay, endDay).
		Group("category").
		Scan(&result).Error
	if err != nil {
		return nil, fmt.Errorf("error getting transactions: %v", err)
	}

	categorySummary := make(map[string]uint64)
	for _, item := range result {
		categorySummary[item.Category] += item.Value
	}
	return categorySummary, nil
}

func GenerateWeeklyExpensesReport(categorySummary map[string]uint64) string {
	if len(categorySummary) == 0 {
		return "📊 За прошедшую неделю расходы отсутствуют."
	}

	report := "📊 Отчёт за неделю:\n\n"
	totalExpense := uint64(0)

	for category, total := range categorySummary {
		report += fmt.Sprintf("▪ Категория: %s — Расход: %d\n", category, total)
		totalExpense += total
	}

	report += fmt.Sprintf("\n💸 Общий расход за неделю составил: %d", totalExpense)
	return report
}

// расход за месяц
func (exp *ExpensesHandler) ExpenseMonthAnalytic(update tgbotapi.Update) (map[string]uint64, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second) // Конец месяца

	var results []struct {
		Category string
		Value    uint64
	}

	err := exp.DB.Model(&models.Transactions{}).
		Select("category, SUM (quantities) as value").
		Where("telegram_id = ? AND operation_type = ? AND created_at BETWEEN ? AND ?",
			update.Message.Chat.ID, false, startOfMonth, endOfMonth).
		Group("category").
		Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("error getting transactions: %v", err)
	}

	categorySummary := make(map[string]uint64)
	for _, item := range results {
		categorySummary[item.Category] += item.Value
	}
	return categorySummary, nil
}

func GenerateMonthlyExpensesReport(categorySummary map[string]uint64) string {
	if len(categorySummary) == 0 {
		return "📊 За прошедший месяц расходы отсутствуют."
	}

	report := "📊 Отчёт за месяц:\n\n"
	totalIncome := uint64(0)

	for category, total := range categorySummary {
		report += fmt.Sprintf("▪ Категория: %s — Расход: %d\n", category, total)
		totalIncome += total
	}

	report += fmt.Sprintf("\n💸 Общий расход за месяц составил: %d", totalIncome)
	return report
}
