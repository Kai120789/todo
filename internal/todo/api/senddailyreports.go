package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"todo/internal/todo/config"
	"todo/internal/todo/models"
	"todo/internal/todo/utils"

	"go.uber.org/zap"
)

func SendDailyReports(tasks []models.Task, chatID int64, status int) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	client := http.Client{}

	type Dto struct {
		ChatID       int64  `json:"chat_id"`
		Message      string `json:"message"`
		MessageEnded string `json:"message_ended"`
	}

	urlString := fmt.Sprintf("%s/scheduler", cfg.TelegramAppURL)

	var message, messageEnded string

	if status == 1 {
		message = utils.FormatTasksMessage(tasks)
	}

	if status == 2 {
		messageEnded = utils.FormatEndedTasksMessage(tasks)
	}

	dto := Dto{
		ChatID:       chatID,
		Message:      message,
		MessageEnded: messageEnded,
	}

	jsonStr, err := json.Marshal(dto)
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
