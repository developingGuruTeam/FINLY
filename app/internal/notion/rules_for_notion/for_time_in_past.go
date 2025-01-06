package rules_for_notion

import (
	"errors"
	"time"
)

func ValidateRightTime(timeStr string) (time.Time, error) {
	date, err := time.Parse("02.01.2006", timeStr)
	if err != nil {
		return time.Time{}, errors.New("Неверный формат даты. Пожалуйста, используйте формат ДД.ММ.ГГГГ")
	}

	now := time.Now()
	if date.Before(now.Truncate(24 * time.Hour)) {
		return time.Time{}, errors.New("Напоминание не может быть в прошлом")
	}

	nextReminder := time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		12, 0, 0, 0,
		time.Local,
	)

	return nextReminder, nil
}
