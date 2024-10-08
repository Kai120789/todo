package api

import (
	"fmt"
	"net/http"
	"net/url"
	"todo/internal/todo/config"
	"todo/internal/todo/models"
	"todo/internal/todo/utils"

	"go.uber.org/zap"
)

// SendDailyReports отправляет ежедневные отчеты в Telegram
func SendDailyReports(tasks []models.Task, chatID int64, status int) error {
	cfg, err := config.GetConfig()
	if err != nil {
		zap.S().Fatal("Error loading config", zap.Error(err))
	}

	// Формирование сообщения в зависимости от статуса задач
	var message string
	if status == 1 {
		message = utils.FormatTasksMessage(tasks)
	} else {
		message = utils.FormatEndedTasksMessage(tasks)
	}

	client := &http.Client{}
	botUrl := fmt.Sprintf("%s/scheduler", cfg.TelegramAppURL)
	taskData := url.Values{
		"chatID": {fmt.Sprintf("%d", chatID)},
		"task":   {message},
	}
	response, err := client.PostForm(botUrl, taskData)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Проверка статуса ответа
	if response.StatusCode == http.StatusOK {
		zap.S().Info("Сообщение успешно отправлено в Telegram!")
	} else {
		zap.S().Errorf("Не удалось отправить сообщение, статус: %s", response.Status)
	}

	return nil
}
