package main

import (
	"cachManagerApp/app/config"
	"cachManagerApp/app/pkg/TgBot"
	"cachManagerApp/app/pkg/logger/slog_init"
	"cachManagerApp/database"
	"sync"
)

func main() {
	log := slog_init.Init()
	log.Info("Logger запущен")
	log.Debug("Цвет прикольный")
	config.LoadEnviroment()
	var wg sync.WaitGroup

	wg.Add(1)

	// устанавливаем соединение с телеграм
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				log.Info("Возникла паника при обработке пользователя: %v", r)
			}
		}()

		if _, err := tg_bot.ConnectToTgBot(log); err != nil {
			log.Error("Ошибка подключения к Telegram боту: %s", "err", err)
		}
	}()

	database.ConnectionDB(log)
	log.Info("БД запущена")

	wg.Wait()
}
