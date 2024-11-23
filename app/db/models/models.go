package models

import (
	"gorm.io/gorm"
	"time"
)

type Transactions struct {
	*gorm.DB
	ID            uint64
	TelegramID    uint64
	CreatedAt     time.Time
	OperationType bool
	Quantities    uint64
	Category      string
	Description   string
}

type Users struct {
	*gorm.DB
	TelegramID   uint64
	Name         string
	Transactions []Transactions
	Сurrency     string
}

type Summary struct {
	Startate     time.tine
	EndDate      time.time
	TotalIncowe  float64
	TotalExpense float64
	Profiti      int64
}
