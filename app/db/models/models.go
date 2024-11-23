package models

import (
	"gorm.io/gorm"
	"time"
)

// тип транзакций
type Transactions struct {
	gorm.Model
	ID            uint64     `gorm:"primaryKey;autoIncrement:true"`
	TelegramID    uint64     `gorm:"foreignKey:telegramID"` // id пользователя
	CreatedAt     time.Time  `gorm:"autoCreateTime"`        // дата/время создания
	OperationType bool       // тип транзакции ( true - доход, false - расход)
	Quantities    uint64     // количество (в валюте)
	Category      Categories `gorm:"not null"` // категория транзакции
	Description   string     // описание транзакции
}

type Users struct {
	gorm.Model
	TelegramID   uint64
	Name         string
	Transactions []Transactions
	// валюта (RUB, USD, EUR, etc.)
	Сurrency string
}

type Categories struct {
	gorm.Model
	Category string `gorm:"uniqueIndex"`
}

//type Summary struct {
//	Startate     time.Time
//	EndDate      time.Time
//	TotalIncowe  float64
//	TotalExpense float64
//	Profiti      int64
//}
