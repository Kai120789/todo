package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"todo/internal/todo/models"

	"go.uber.org/zap"
)

func Create(task models.Task, chatID int64) error {

	type TaskDtoChatID struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		StatusId    uint   `json:"status_id"`
		ChatId      int64
	}

	client := &http.Client{}
	createURL := fmt.Sprintf("%s/create-task", "http://host.docker.internalt:8081")

	dto := TaskDtoChatID{
		Title:       task.Title,
		Description: task.Description,
		StatusId:    task.StatusId,
		ChatId:      chatID,
	}

	jsonStr, err := json.Marshal(dto)
	if err != nil {
		zap.S().Error("error marshalling DTO", zap.Error(err))
		return err
	}

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
