package models

import (
	"time"
)

// тип транзакций
type Transactions struct {
	ID            uint64
	TelegramID    uint64
	CreatedAt     time.Time
	OperationType bool
	Quantities    uint64
	CategoryID    uint
	Category      Categories
	Description   string
}

type Users struct {
	TelegramID   uint64
	Name         string
	Transactions []Transactions
	Сurrency     string
}

type Categories struct {
	Category string
}

//type Summary struct {
//	Startate     time.Time
//	EndDate      time.Time
//	TotalIncowe  float64
//	TotalExpense float64
//	Profiti      int64
//}
