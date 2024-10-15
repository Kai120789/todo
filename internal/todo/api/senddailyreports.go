package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"todo/internal/todo/config"
	"todo/internal/todo/models"

	"go.uber.org/zap"
)

type MessDto struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	StatusId    uint   `json:"status_id"`
	ChatId      int64  `json:"chat_id"`
}

func SendDailyReports(tasks []models.Task, chatID int64, status int) error {
	client := http.Client{}

	var messDto []MessDto

	urlString := fmt.Sprintf("%s/scheduler", config.AppConfig.TelegramAppURL)

	for _, task := range tasks {
		task := MessDto{
			Title:       task.Title,
			Description: task.Description,
			StatusId:    task.StatusId,
			ChatId:      chatID,
		}
		messDto = append(messDto, task)
	}

	jsonStr, err := json.Marshal(messDto)
	if err != nil {
		zap.S().Error("error marshalling DTO", zap.Error(err))
		return err
	}

	resp, err := client.Post(urlString, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		zap.S().Error("error during user registration", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.S().Error("error reading response body", zap.Error(err))
	}
	fmt.Println(string(body))

	return nil

}
