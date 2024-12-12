package methodsForIncomeAnalys

import (
	"cachManagerApp/app/db/models"
	"cachManagerApp/app/pkg/logger"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"time"
)

type IncomeAnalyticHandlers interface {
	IncomeDayAnalytic(update tgbotapi.Update) ([]models.Transactions, error)
	IncomeWeekAnalytic(update tgbotapi.Update) (map[string]uint64, error)
	IncomeMonthAnalytic(update tgbotapi.Update) (map[string]uint64, error)
}

var log = logger.GetLogger()

type AnalyticHandler struct {
	DB *gorm.DB
}

func (an *AnalyticHandler) IncomeDayAnalytic(update tgbotapi.Update) ([]models.Transactions, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Запрос к базе данных
	var transactions []models.Transactions
	err := an.DB.Where("telegram_id = ? AND operation_type = ? AND created_at BETWEEN ? AND ?",
		update.Message.Chat.ID, true, startOfDay, endOfDay).Find(&transactions).Error
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении транзакций: %w", err)
	}
	return transactions, nil
}

func GenerateDailyIncomeReport(transactions []models.Transactions) string {
	if len(transactions) == 0 {
		return "📈 Сегодня у вас не было доходов."
	}

	report := "📈 Отчёт за день:\n\n"
	totalIncome := uint64(0)

	for _, t := range transactions {
		report += fmt.Sprintf("▪ Категория: %s\n", t.Category)
		report += fmt.Sprintf("   Сумма: %d\n", t.Quantities)
		if t.Description != "" {
			report += fmt.Sprintf("   Комментарий: %s\n", t.Description)
		}
		report += "\n"
		totalIncome += t.Quantities
	}

	report += fmt.Sprintf("💵 Итого доходов за день: %d\n", totalIncome)
	return report
}

func (an *AnalyticHandler) IncomeWeekAnalytic(update tgbotapi.Update) (map[string]uint64, error) {
	now := time.Now()
	startOfWeek := now.AddDate(0, 0, -7)
	endOfWeek := now

	var results []struct {
		Category   string
		TotalValue uint64
	}

	err := an.DB.Model(&models.Transactions{}).
		Select("category, SUM (quantities) as total_value").
		Where("telegram_id = ? AND operation_type = ? AND created_at BETWEEN ? AND ?",
			update.Message.Chat.ID, true, startOfWeek, endOfWeek).
		Group("category").
		Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении данных за неделю: %w", err)
	}

	categorySummary := make(map[string]uint64)
	for _, result := range results {
		categorySummary[result.Category] += result.TotalValue
	}
	return categorySummary, nil
}

func GenerateWeeklyIncomeReport(categorySummary map[string]uint64) string {
	if len(categorySummary) == 0 {
		return "📊 За прошедшую неделю доходы отсутствуют."
	}

	report := "📊 Отчёт за неделю:\n\n"
	totalIncome := uint64(0)

	for category, total := range categorySummary {
		report += fmt.Sprintf("▪ Категория: %s — Доход: %d\n", category, total)
		totalIncome += total
	}

	report += fmt.Sprintf("\n💵 Общий доход за неделю составил: %d", totalIncome)
	return report
}

func (an *AnalyticHandler) IncomeMonthAnalytic(update tgbotapi.Update) (map[string]uint64, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second) // Конец месяца

	var results []struct {
		Category   string
		TotalValue uint64
	}

	err := an.DB.Model(&models.Transactions{}).
		Select("category, SUM (quantities) as total_value").
		Where("telegram_id = ? AND operation_type = ? AND created_at BETWEEN ? AND ?",
			update.Message.Chat.ID, true, startOfMonth, endOfMonth).
		Group("category").
		Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении данных за неделю: %w", err)
	}

	categorySummary := make(map[string]uint64)
	for _, result := range results {
		categorySummary[result.Category] += result.TotalValue
	}
	return categorySummary, nil
}

func GenerateMonthlyIncomeReport(categorySummary map[string]uint64) string {
	if len(categorySummary) == 0 {
		return "📊 За прошедший месяц доходы отсутствуют."
	}

	report := "📊 Отчёт за месяц:\n\n"
	totalIncome := uint64(0)

	for category, total := range categorySummary {
		report += fmt.Sprintf("▪ Категория: %s — Доход: %d\n", category, total)
		totalIncome += total
	}

	report += fmt.Sprintf("\n💵 Общий доход за месяц составил: %d", totalIncome)
	return report
}
