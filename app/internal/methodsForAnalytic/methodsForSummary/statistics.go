package methodsForSummary

import (
	"cachManagerApp/app/db/models"
	"fmt"
	"gorm.io/gorm"
)

func GenerateStatisticsReport(userID int64, db *gorm.DB) string {
	report := "🧮 *Статистика*\n\n"

	var user models.Users
	db.First(&user, "telegram_id = ?", userID)

	// имя
	name := user.Name

	// дата регистрации
	registrationDate := user.CreatedAt.Format("02-01-2006")

	// валюта
	currency := user.Currency

	// оборот
	var allIncomes, allExpenses int64
	db.Model(&models.Transactions{}).
		Where("telegram_id = ? AND operation_type = ?", userID, true).
		Select("SUM(quantities)").
		Scan(&allIncomes)
	db.Model(&models.Transactions{}).
		Where("telegram_id = ? AND operation_type = ?", userID, false).
		Select("SUM(quantities)").
		Scan(&allExpenses)
	allBalance := allIncomes - allExpenses

	// макс доход
	var maxIncome uint64
	db.Table("transactions").
		Select("MAX(quantities)").Where("telegram_id = ? AND operation_type = ?", userID, true).
		Scan(&maxIncome)

	// макс расход
	var maxExpense uint64
	db.Table("transactions").
		Select("MAX(quantities)").Where("telegram_id = ? AND operation_type = ?", userID, false).
		Scan(&maxExpense)

	// топ категория доходов
	var categoryInc string
	var totalInc uint64
	db.Table("transactions").
		Select("category, SUM(quantities) as Total").
		Where("telegram_id = ? AND operation_type = ?", userID, true).
		Group("category").
		Order("Total DESC").
		Limit(1).
		Row().Scan(&categoryInc, &totalInc)

	// топ категория доходов
	var categoryExp string
	var totalExp uint64
	db.Table("transactions").
		Select("category, SUM(quantities) as Total").
		Where("telegram_id = ? AND operation_type = ?", userID, false).
		Group("category").
		Order("Total DESC").
		Limit(1).
		Row().Scan(&categoryExp, &totalExp)

	// кол-во операций всего в штуках
	var incCount, expCount int64
	db.Model(&models.Transactions{}).
		Where("telegram_id = ? AND operation_type = ?", userID, true).
		Count(&incCount)
	db.Model(&models.Transactions{}).
		Where("telegram_id = ? AND operation_type = ?", userID, false).
		Count(&expCount)

	report += fmt.Sprintf("👤 Имя: *%s*\n\n", name)
	report += fmt.Sprintf("📅 Дата регистрации: *%s*\n\n", registrationDate)
	report += fmt.Sprintf("💱 Текущая валюта: *%s*\n\n", currency)
	report += fmt.Sprintf("⚖️ Баланс за все время *%d* %s\n\n", allBalance, currency)
	report += fmt.Sprintf("🟢️ Максимальный доход: *%d* %s\n\n", maxIncome, currency)
	report += fmt.Sprintf("🔴 Максимальный расход: *%d* %s\n\n", maxExpense, currency)
	report += fmt.Sprintf("📥 Основной доходы в категории: *%s %d* %s\n\n", categoryInc, totalInc, currency)
	report += fmt.Sprintf("📤 Основые расходы в категории: *%s %d* %s\n\n", categoryExp, totalExp, currency)
	report += fmt.Sprintf("🗃 Всего операций: *%v* шт. \n(Доходы: *%v* шт. Расходы: *%v* шт.)", expCount+incCount, incCount, expCount)

	return report
}
