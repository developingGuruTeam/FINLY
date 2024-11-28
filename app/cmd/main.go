package main

import (
	"cachManagerApp/app/pkg/TgBot"
	"cachManagerApp/database"
)

func main() {
	TgBot.ConnectToTgBot()

	database.Connect()

	// Создание таблиц в базе данных
	database.AutoMigrate()
}
