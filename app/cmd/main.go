package main

import (
	"cachManagerApp/app/pkg/TgBot"
	"cachManagerApp/database"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}

	go func() {
		if _, err := TgBot.ConnectToTgBot(); err != nil {
			log.Fatalf("Ошибка подключения к Telegram боту: %v", err)
		}
	}()

	database.ConnectionDB()
	log.Println("ДБ запущена")
	// database.AutoMigrate()
}
