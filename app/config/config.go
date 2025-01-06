package config

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnviroment() {
	// загружаем переменные окружения
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}
}
