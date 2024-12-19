package models

import (
	"time"
)

// тип транзакций
type Transactions struct {
	ID            uint64    `gorm:"primary_key;auto_increment" json:"id"`
	TelegramID    uint64    `gorm:"not null" json:"telegram_id"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	OperationType bool      `gorm:"not null" json:"operation_type"`
	Quantities    uint64    `gorm:"not null" json:"quantities"`
	Category      string    `gorm:"not null" json:"category_id"`
	Description   string    `gorm:"type:text" json:"description"`
}

// структура пользователя
type Users struct {
	TelegramID uint64 `gorm:"primary_key" json:"telegram_id"`
	Name       string `gorm:"not null" json:"name"`
	Currency   string `gorm:"not null" json:"currency"`
}

// структура категории
type CategorySummary struct {
	Category string
	Amount   uint64
}

// структура для сальдо
type Summary struct {
	TotalIncome       uint64
	TotalExpense      uint64
	Profit            int64
	TopIncome         CategorySummary
	TopExpense        CategorySummary
	IncomeCategories  []CategorySummary
	ExpenseCategories []CategorySummary
}

// структура напоминания
type Reminder struct {
	ID            uint      `gorm:"primary_key" json:"id"`
	UserID        uint64    `gorm:"not null" json:"user_id"`          // Telegram ID пользователя
	Amount        int       `gorm:"not null" json:"amount"`           // Сумма платежа
	Frequency     string    `gorm:"not null" json:"frequency"`        // Периодичность (неделя/месяц)
	NextReminder  time.Time `gorm:"not null" json:"next_reminder"`    // Дата следующего напоминания
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"` // Дата создания
	LastUpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"` // Последнее обновление
	Category      string    `gorm:"not null" json:"category"`
}
