package methodsForSummary

import (
	"cachManagerApp/app/db/models"
	"fmt"
	"gorm.io/gorm"
)

func GenerateStatisticsReport(userID int64, db *gorm.DB) string {
	report := "🧮 *Статистика*\n"

	/*
		1. Имя +
		2. Дата регистрации +
		3. Текущая валюта +
		4. Валовый оборот все приходы все расходы и баланс на сейчас за все время
		5. Максимальный доход
		6. Максимальный расход
		7. Самый прибыльный месяц
		8. Самый затратный месяц
		9. Основная категория доходов
		10. Любимые расходы

	*/

	// создание юзера
	var user models.Users
	db.First(&user, "telegram_id = ?", userID)

	// имя
	name := user.Name

	// дата регистрации
	registrationDate := user.CreatedAt.Format("02-01-2006")

	// валюта
	currency := user.Currency

	report += fmt.Sprintf("👤 Имя: *%s*\n", name)
	report += fmt.Sprintf("📅 Дата регистрации: *%s*\n", registrationDate)
	report += fmt.Sprintf("💱 Текущая валюта: *%s*\n", currency)

	return report
}
