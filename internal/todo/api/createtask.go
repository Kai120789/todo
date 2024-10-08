package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"todo/internal/todo/config"
	"todo/internal/todo/models"

	"go.uber.org/zap"
)

func Create(task *models.Task, chatID int64) error {
	// Получаем конфигурацию
	cfg, err := config.GetConfig()
	if err != nil {
		zap.S().Fatalf("get config error", zap.Error(err))
	}

	// Формируем данные для отправки в JSON формате
	taskData := map[string]interface{}{
		"chatID":  chatID,
		"task":    task.Title,
		"details": task.Description,
		"status":  task.StatusId,
	}

	jsonData, err := json.Marshal(taskData)
	if err != nil {
		zap.S().Error("failed to marshal task data", zap.Error(err))
		return err
	}

	// Отправляем POST запрос
	client := &http.Client{}
	botUrl := fmt.Sprintf("%s/create-task", cfg.TelegramAppURL)
	request, err := http.NewRequest("POST", botUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		zap.S().Error("failed to create POST request", zap.Error(err))
		return err
	}

	// Устанавливаем заголовки
	request.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	response, err := client.Do(request)
	if err != nil {
		zap.S().Error("failed to send POST request", zap.Error(err))
		return err
	}
	defer response.Body.Close()

	// Проверка статуса ответа
	if response.StatusCode != http.StatusOK {
		zap.S().Errorf("received non-OK status: %d", response.StatusCode)
		return fmt.Errorf("received non-OK status: %d", response.StatusCode)
	}

	return nil
}
