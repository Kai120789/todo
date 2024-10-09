package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"todo/internal/todo/config"
	"todo/internal/todo/models"

	"go.uber.org/zap"
)

func Create(task models.Task, chatID int64) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	client := &http.Client{}
	createURL := fmt.Sprintf("%s/create-task", cfg.TelegramAppURL)

	var jsonStr = []byte(fmt.Sprintf(`{"title":"%s", "description":"%s", "status_id":%d}`, task.Title, task.Description, task.StatusId))

	// Создание io.Reader из JSON
	response, err := client.Post(createURL, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		zap.S().Error("error during user registration", zap.Error(err))
		return err
	}
	defer response.Body.Close()

	fmt.Println("Response status:", response.StatusCode)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		zap.S().Error("error reading response body", zap.Error(err))
	}
	fmt.Println(string(body))

	return err
}
