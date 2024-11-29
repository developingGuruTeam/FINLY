package database

import (
	"cachManagerApp/app/db/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

func ConnectionDB() {

	fmt.Printf("DB_HOST=%s, DB_USER=%s,DB_PASSWORD=%s, DB_NAME=%s, DB_PORT=%s\n",
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
		log.Fatal("Не удалось подключиться к БД", err)
		return
	}

	DB = db
	fmt.Println("Подключение к БД успешно")

	err = db.AutoMigrate(&models.Users{}, &models.Categories{}, &models.Transactions{})
	if err != nil {
		log.Fatal("Ошибка при миграции", err)
		return
	}

	fmt.Println("Миграции успешно выполнены")
}
