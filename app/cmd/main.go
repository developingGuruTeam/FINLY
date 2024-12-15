package main

import (
	"cachManagerApp/app/pkg/TgBot"
	"cachManagerApp/database"
	"github.com/joho/godotenv"
	"log"
	"sync"
)

func main() {
	// загружаем переменные окружения
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}
	var wg sync.WaitGroup

	wg.Add(1)

	// устанавливаем соединение с телеграм
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Возникла паника при обработке пользователя: %v", r)
			}
		}()

		if _, err := TgBot.ConnectToTgBot(); err != nil {
			log.Fatalf("Ошибка подключения к Telegram боту: %v", err)
		}
	}()
	database.ConnectionDB()
	log.Println("БД запущена")

	wg.Wait()

}
