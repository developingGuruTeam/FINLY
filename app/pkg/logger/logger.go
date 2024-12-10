package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	// Инициализация логгера
	log = logrus.New()

	// Открытие файла для записи логов
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Ошибка при открытии файла: %v", err)
	}

	// Настройка логгера
	log.SetOutput(file)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(logrus.InfoLevel)
}

// GetLogger возвращает singleton логгера
func GetLogger() *logrus.Logger {
	return log
}
