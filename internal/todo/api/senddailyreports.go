package api

import (
	"todo/internal/todo/models"
)

// SendDailyReports отправляет ежедневные отчеты в Telegram
func SendDailyReports(tasks []models.Task, chatID int64, status int) error {
	return nil
}
