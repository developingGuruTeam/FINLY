package models

import (
	"gorm.io/gorm"
	"time"
)

// тип транзакций
type Transactions struct {
	gorm.Model
	ID            uint64     `gorm:"primaryKey;autoIncrement:true"`
	TelegramID    uint64     `gorm:"not null"`       // Внешний ключ на Users.TelegramID
	CreatedAt     time.Time  `gorm:"autoCreateTime"` // Дата/время создания
	OperationType bool       // Тип транзакции (доход/расход)
	Quantities    uint64     // Количество (в валюте)
	CategoryID    uint       `gorm:"not null"`              // Внешний ключ на Categories
	Category      Categories `gorm:"foreignKey:CategoryID"` // Связь с таблицей Categories
	Description   string     // Описание транзакции
}

type Users struct {
	gorm.Model
	TelegramID   uint64         `gorm:"uniqueIndex"` // Уникальный TelegramID для связи
	Name         string         // Имя пользователя
	Transactions []Transactions `gorm:"foreignKey:TelegramID;references:TelegramID"` // Указание внешнего ключа
	Сurrency     string         // Валюта (RUB, USD, EUR)
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
