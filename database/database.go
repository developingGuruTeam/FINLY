package database

import (
	"cachManagerApp/app/db/models"
	"fmt"
	"log/slog"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectionDB(log *slog.Logger) {

	log.Info("DB_HOST=%s, DB_USER=%s,DB_PASSWORD=%s, DB_NAME=%s, DB_PORT=%s",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error("Не удалось подключиться к БД", log.With("error", err))
		return
	}

	DB = db
	log.Info("Подключение к БД успешно")

	err = db.AutoMigrate(&models.Users{}, &models.Transactions{})
	if err != nil {
		log.Error("Ошибка при миграции", log.With("error", err))
		return
	}

	err = db.AutoMigrate(&models.Reminder{})
	if err != nil {
		log.Error("Ошибка при миграции", log.With("error", err))
		return
	}

	log.Info("Миграции успешно выполнены")

	// Очистка таблиц

	// Очистка операций
	//db.Exec("DELETE FROM transactions")

	// Очистка таблицы пользователей
	//db.Exec("DELETE FROM users")

}
