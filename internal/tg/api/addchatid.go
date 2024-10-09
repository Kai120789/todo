package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

func AddChatID(username string, chatID int64, appURL string) bool {
	client := &http.Client{}
	registerURL := fmt.Sprintf("%s/addchatid", appURL)

	var jsonStr = []byte(fmt.Sprintf(`{"tg_name":"%s", "chat_id":%d}`, username, chatID))

	// Создание io.Reader из JSON
	response, err := client.Post(registerURL, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		zap.S().Error("error during user registration", zap.Error(err))
		return false
	}
	defer response.Body.Close()

	fmt.Println("Response status:", response.StatusCode)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		zap.S().Error("error reading response body", zap.Error(err))
	}
	fmt.Println(string(body))

	return response.StatusCode == http.StatusCreated
}
