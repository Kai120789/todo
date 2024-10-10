package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"todo/internal/tg/dto"

	"go.uber.org/zap"
)

func AddChatID(username string, chatID int64, appURL string) error {
	client := &http.Client{}
	registerURL := fmt.Sprintf("%s/add-chat-id", appURL)

	dto := dto.ChatID{
		Username: username,
		ChatID:   chatID,
	}

	jsonStr, err := json.Marshal(dto)
	if err != nil {
		zap.S().Error("error marshalling DTO", zap.Error(err))
		return err
	}

	// Создание io.Reader из JSON
	response, err := client.Post(registerURL, "application/json", bytes.NewBuffer(jsonStr))
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

	if response.StatusCode != http.StatusCreated {
		return err
	}
	return nil
}
