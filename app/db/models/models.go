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
	CategoryID    uint      `gorm:"not null" json:"category_id"`
	Description   string    `gorm:"type:text" json:"description"`
}

type Users struct {
	TelegramID uint64 `gorm:"primary_key" json:"telegram_id"`
	Name       string `gorm:"not null" json:"name"`
	Currency   string `gorm:"not null" json:"currency"`
}

//  type Summary struct {
//	Startate     time.Time
//	EndDate      time.Time
//	TotalIncowe  float64
//	TotalExpense float64
//	Profit       int64
//}
