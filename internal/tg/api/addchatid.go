package api

import (
	"fmt"
	"net/http"
	"net/url"

	"go.uber.org/zap"
)

func AddChatID(username string, chatID int64, appURL string) bool {
	client := &http.Client{}
	registerURL := fmt.Sprintf("%s/addchatid", appURL)

	data := url.Values{
		"username": {username},
		"chatID":   {fmt.Sprintf("%d", chatID)},
	}

	response, err := client.PostForm(registerURL, data)
	if err != nil {
		zap.S().Error("error during user registration", zap.Error(err))
		return false
	}
	defer response.Body.Close()

	return response.StatusCode == http.StatusOK
}
