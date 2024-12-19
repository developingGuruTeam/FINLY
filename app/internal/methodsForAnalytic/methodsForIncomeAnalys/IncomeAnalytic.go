package methodsForIncomeAnalys

import (
	"cachManagerApp/app/db/models"
	"cachManagerApp/app/pkg/logger"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"time"
)

// хендлер доходов
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

func GenerateDailyIncomeReport(transactions []models.Transactions, currency string) string {
	if len(transactions) == 0 {
		return "📈 За сегодня доходов нет"
	}

	report := "📈 *Отчёт за день (доходы)*\n\n"
	totalIncome := uint64(0)

	for _, inc := range transactions {
		// +3 часа к времени чтобы было по мск делаю временно. НАДО В БД ПОСТАВИТЬ НАШЕ ВРЕМЯ по умолчанию!!! либо как-то к юсеру привязаться
		localTime := inc.CreatedAt.Add(3 * time.Hour)
		formattedTime := localTime.Format("15:04")

		report += fmt.Sprintf("▪ *%s*\n", inc.Category)
		report += fmt.Sprintf("%d %s 📝 _%v_", inc.Quantities, currency, formattedTime)

		// сокращаем коммент пользователя на вывод
		if inc.Description != "" {
			decs := inc.Description
			runes := []rune(decs)
			if len([]rune(decs)) > 32 {
				decs = string(runes[:32])
			}

			report += fmt.Sprintf(" _%s_", decs)
		}
		report += "\n"
		totalIncome += inc.Quantities
	}

	report += fmt.Sprintf("\n💸 Итого расходов за день:\n*%d* %s\n", totalIncome, currency)
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

func GenerateWeeklyIncomeReport(categorySummary map[string]uint64, currency string) string {
	categoryDetails := map[string]string{
		"Заработная плата":    "🔵",
		"Побочный доход":      "🔴",
		"Доход от бизнеса":    "🟡",
		"Гос. выплаты":        "🟢",
		"Продажа имущества":   "🟠",
		"Доход от инвестиций": "🟣",
		"Прочие доходы":       "⚪️",
	}

	if len(categorySummary) == 0 {
		return "📊 За прошедшую неделю доходы отсутствуют."
	}

	totalIncome := uint64(0)
	for _, value := range categorySummary {
		totalIncome += value
	}

	report := "📊 *Доходы за неделю*\n\n"

	for category, value := range categorySummary {
		percentage := (float64(value) / float64(totalIncome)) * 100
		if emoji, exists := categoryDetails[category]; exists {
			report += fmt.Sprintf("%s %s: %d %s (%d%%)\n", emoji, category, value, currency, int(percentage))
		} else {
			report += fmt.Sprintf("%s: %d %s (%d%%)\n", category, value, currency, int(percentage))
		}
	}

	report += fmt.Sprintf("\n💸 Общий доход за неделю: *%d* %s", totalIncome, currency)
	return report
}

func (a *AnalyticHandler) IncomeMonthAnalytic(update tgbotapi.Update) (map[string]uint64, uint64, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second)

	var results []struct {
		Category string
		Value    uint64
	}

	err := a.DB.Model(&models.Transactions{}).
		Select("category, SUM(quantities) as value").
		Where("telegram_id = ? AND operation_type = ? AND created_at BETWEEN ? AND ?",
			update.Message.Chat.ID, true, startOfMonth, endOfMonth). // Только доходы
		Group("category").
		Scan(&results).Error

	if err != nil {
		return nil, 0, fmt.Errorf("ошибка анализа доходов за месяц: %v", err)
	}

	categorySummary := make(map[string]uint64)
	totalIncome := uint64(0)

	for _, item := range results {
		categorySummary[item.Category] = item.Value
		totalIncome += item.Value
	}

	return categorySummary, totalIncome, nil
}

func GenerateMonthlyIncomeReport(categorySummary map[string]uint64, currency string) string {
	categoryDetails := map[string]string{
		"Заработная плата":    "🔵",
		"Побочный доход":      "🔴",
		"Доход от бизнеса":    "🟡",
		"Гос. выплаты":        "🟢",
		"Продажа имущества":   "🟠",
		"Доход от инвестиций": "🟣",
		"Прочие доходы":       "⚪️",
	}

	if len(categorySummary) == 0 {
		return "📊 За прошедший месяц доходы отсутствуют."
	}

	totalIncome := uint64(0)
	for _, value := range categorySummary {
		totalIncome += value
	}

	report := "📊 *Доходы за месяц*\n\n"

	for category, value := range categorySummary {
		percentage := (float64(value) / float64(totalIncome)) * 100
		if emoji, exists := categoryDetails[category]; exists {
			report += fmt.Sprintf("%s %s: %d (%d%%)\n", emoji, category, value, int(percentage))
		} else {
			report += fmt.Sprintf("%s: %d (%d%%)\n", category, value, int(percentage))
		}
	}

	report += fmt.Sprintf("\n💸 Общий доход за месяц: *%d* %s", totalIncome, currency)

	return report
}
