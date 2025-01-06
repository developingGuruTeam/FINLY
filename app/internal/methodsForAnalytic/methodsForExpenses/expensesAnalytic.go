package methodsForExpenses

import (
	"cachManagerApp/app/db/models"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

type ExpensesHandler struct {
	DB *gorm.DB
}

// хендлер расходов
//
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
		return nil, fmt.Errorf("ошибка заполнения модели транзакции в обработке аналитики по расходам за день: %v", err)
	}

	return transactions, nil
}

func GenerateDailyExpenseReport(expenses []models.Transactions, currency string) string {
	if len(expenses) == 0 {
		return "📉 За сегодня расходов нет."
	}

	report := "📉 *Отчёт за день (расходы)*\n\n"
	var totalExpenses uint64

	for _, exp := range expenses {
		// +3 часа к времени чтобы было по мск делаю временно. НАДО В БД ПОСТАВИТЬ НАШЕ ВРЕМЯ по умолчанию!!! либо как-то к юсеру привязаться
		localTime := exp.CreatedAt.Add(3 * time.Hour)
		formattedTime := localTime.Format("15:04")

		report += fmt.Sprintf("▪ *%s*\n", exp.Category)
		report += fmt.Sprintf("%d %s 📝 _%v_", exp.Quantities, currency, formattedTime)
		// сокращаем коммент пользователя на вывод
		if exp.Description != "" {
			decs := exp.Description
			runes := []rune(decs)
			if len([]rune(decs)) > 32 {
				decs = string(runes[:32])
			}

			report += fmt.Sprintf(" _%s_", decs)
		}
		report += "\n"
		totalExpenses += exp.Quantities
	}
	report += fmt.Sprintf("\n💸 Итого расходов за день:\n*%d* %s\n", totalExpenses, currency)
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
		return nil, fmt.Errorf("ошибка в заполнении транзакции по расходам за неделю: %v", err)
	}

	categorySummary := make(map[string]uint64)
	for _, item := range result {
		categorySummary[item.Category] += item.Value
	}
	return categorySummary, nil
}

func GenerateWeeklyExpensesReport(categorySummary map[string]uint64, currency string) string {
	categoryDetails := map[string]string{
		"Бытовые траты":       "🔵",
		"Регулярные платежи":  "🔴",
		"Одежда":              "🟡",
		"Здоровье":            "🟢",
		"Досуг и образование": "🟠",
		"Инвестиции":          "🟣",
		"Прочие расходы":      "⚪️",
	}

	if len(categorySummary) == 0 {
		return "📊 За прошедшую неделю расходы отсутствуют."
	}

	totalExpense := uint64(0)
	for _, value := range categorySummary {
		totalExpense += value
	}

	report := "📊 *Расходы за неделю*\n\n"

	for category, value := range categorySummary {
		percentage := (float64(value) / float64(totalExpense)) * 100
		if emoji, exists := categoryDetails[category]; exists {
			report += fmt.Sprintf("%s %s: %d %s (%d%%)\n", emoji, category, value, currency, int(percentage))
		} else {
			report += fmt.Sprintf("%s: %d %s (%d%%)\n", category, value, currency, int(percentage))
		}
	}

	report += fmt.Sprintf("\n💸 Общий расход за неделю: *%d* %s", totalExpense, currency)
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
		Select("category, SUM(quantities) as value").
		Where("telegram_id = ? AND operation_type = ? AND created_at BETWEEN ? AND ?",
			update.Message.Chat.ID, false, startOfMonth, endOfMonth). // Только расходы
		Group("category").
		Scan(&results).Error

	log.Printf("Результаты запроса за месяц: %+v", results) // Логирование

	if err != nil {
		return nil, fmt.Errorf("ошибка по расходам за месяц: %v", err)
	}

	categorySummary := make(map[string]uint64)
	for _, item := range results {
		categorySummary[item.Category] += item.Value
	}
	return categorySummary, nil
}

func GenerateMonthlyExpensesReport(categorySummary map[string]uint64, currency string) string {
	categoryDetails := map[string]string{
		"Бытовые траты":       "🔵",
		"Регулярные платежи":  "🔴",
		"Одежда":              "🟡",
		"Здоровье":            "🟢",
		"Досуг и образование": "🟠",
		"Инвестиции":          "🟣",
		"Прочие расходы":      "⚪️",
	}

	if len(categorySummary) == 0 {
		return "📊 За прошедший месяц расходы отсутствуют."
	}

	// общий расход
	totalExpense := uint64(0)
	for _, value := range categorySummary {
		totalExpense += value
	}

	report := "📊 *Расходы за месяц*\n\n"

	for category, value := range categorySummary {
		// считаем проценты
		percentage := (float64(value) / float64(totalExpense)) * 100

		// Добавляем строку отчёта
		if emoji, exists := categoryDetails[category]; exists {
			report += fmt.Sprintf("%s %s: %d (%d%%)\n", emoji, category, value, int(percentage))
		} else {
			report += fmt.Sprintf("%s: %d (%d%%)\n", category, value, int(percentage))
		}
	}

	// финиш
	report += fmt.Sprintf("\n💸 Общие расходы за месяц: *%d* %s", totalExpense, currency)

	return report
}
