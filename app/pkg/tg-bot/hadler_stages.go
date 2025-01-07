package tg_bot

import (
	"cachManagerApp/app/internal/methodsForTransaction"
	"cachManagerApp/app/pkg/ButtonsCreate"
	"log/slog"
	"strconv"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommentResponse struct {
	Category string `json:"category"`
	Amount   int64  `json:"amount"`
}

type TransactionResponse struct {
	Action string `json:"action"`
}

type UserResponse struct {
	Action string `json:"action"`
}

var (
	mu                sync.Mutex                            // мьютекс для синхронизации доступа к мапе
	commentStates     = make(map[int64]CommentResponse)     // мапа для хранения состояния комментариев
	userStates        = make(map[int64]UserResponse)        // мапа для хранения состояния пользователей
	transactionStates = make(map[int64]TransactionResponse) // мапа для хранения состояния транзакций
)

// обработка нажатий на кнопки (команда приходит сюда)
func PushOnButton(bot *tgbotapi.BotAPI, update tgbotapi.Update, buttonCreator buttons_create.TelegramButtonCreator, log *slog.Logger) {
	if update.Message != nil {
		chatID := update.Message.Chat.ID

		// блокируем доступ к общим мапам для синхронизации
		mu.Lock()
		val, ok := userStates[chatID]          // проверяем, активен ли режим смены имени/валюты
		val2, ok2 := transactionStates[chatID] // проверяем, есть ли состояние транзакции
		val3, ok3 := commentStates[chatID]     // проверяем, ожидается ли ввод комментария
		mu.Unlock()

		// если активен режим ожидания комментария
		if ok3 {
			if update.Message.Text == "⤵️ Пропустить" {
				// сохраняем транзакцию без коммента
				transaction := methods_for_transactions.TransactionsMethod{}
				if err := transaction.PostTransactionWithComment(update, val3.Category, val3.Amount, "", log); err != nil {
					log.Info("Failed to save transaction without comment: %s", slog.Any("error", err))
					msg := tgbotapi.NewMessage(chatID, "❌ Ошибка сохранения транзакции.")
					bot.Send(msg)
					return
				}

				doneMsg := "✅ Сумма сохранена."
				returnToMainMenu(bot, chatID, buttonCreator, doneMsg) // через функцию возвращаем в главное меню
				mu.Lock()
				delete(commentStates, chatID) // удаляем состояние ожидания комментария
				mu.Unlock()
				return
			}

			// сохраняем транзакцию с коммента
			comment := update.Message.Text
			transaction := methods_for_transactions.TransactionsMethod{}
			if err := transaction.PostTransactionWithComment(update, val3.Category, val3.Amount, comment, log); err != nil {
				log.Info("Failed to save transaction with comment: %s", slog.Any("error", err))
				msg := tgbotapi.NewMessage(chatID, "❌ Ошибка сохранения транзакции.")
				bot.Send(msg)
				return
			}

			doneMsg := "✅ Сумма сохранена.\n📝 Комментарий добавлен"
			returnToMainMenu(bot, chatID, buttonCreator, doneMsg)
			mu.Lock()
			delete(commentStates, chatID)
			mu.Unlock()
			return
		}

		// если активна транзакция, но комментарий еще не введен
		if ok2 && val2.Action != "" {
			// проверка на число
			sum, err := strconv.Atoi(update.Message.Text)
			if err != nil || sum <= 0 {
				msg := tgbotapi.NewMessage(chatID, "🚫 Введите корректное положительное целое число.")
				bot.Send(msg)
				return
			}

			// сохраняем сумму и категорию в состояние комментария
			mu.Lock()
			commentStates[chatID] = CommentResponse{
				Category: val2.Action,
				Amount:   int64(sum),
			}
			delete(transactionStates, chatID) // удаляем состояние транзакции, чтобы дальше запросить коммент
			mu.Unlock()

			msg := tgbotapi.NewMessage(chatID, "Добавьте комментарий к сумме или нажмите ⤵️*Пропустить*")
			msg.ParseMode = "Markdown"
			msg.ReplyMarkup = buttonCreator.CreateCommentButtons() // добавляем на экран кнопку пропустить
			bot.Send(msg)
			return
		}

		// если активен режим смены имени или валюты
		if ok && val.Action != "" {
			handleUserAction(bot, update, val, buttonCreator, log) // запуск через отдельную функцию
			return
		}

		// запускаем обработчик нажатия на кнопки
		handleButtonPress(bot, update, buttonCreator, log)
	}
}
