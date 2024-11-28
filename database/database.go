package database

import (
	"cachManagerApp/app/db/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Подключение к базе данных
func Connect() {
	// Получение данных для подключения из переменных окружения
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),    // Адрес хоста
		getEnv("DB_USER", "postgres"),     // Пользователь
		getEnv("DB_PASSWORD", "08092001"), // Пароль
		getEnv("DB_NAME", "cash_manager"), // Имя базы данных
		getEnv("DB_PORT", "5432"),         // Порт
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}

	log.Println("Успешное подключение к базе данных!")
}

// Вспомогательная функция для получения переменных окружения
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func AutoMigrate() {
	err := DB.AutoMigrate(
		&models.Users{},
		&models.Transactions{},
		&models.Categories{},
	)
	if err != nil {
		log.Fatalf("Ошибка при создании таблиц: %v", err)
	}

	log.Println("Таблицы успешно созданы!")
}
