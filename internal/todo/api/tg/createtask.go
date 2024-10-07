package tg

import (
	"fmt"
	"net/http"
	"net/url"
	"todo/internal/tg/config"
	"todo/internal/todo/models"

	"go.uber.org/zap"
)

func Create(task *models.Task, chatID int64) error {
	// init config
	cfg, err := config.GetConfig()
	if err != nil {
		zap.S().Fatal("error load config", zap.Error(err))
	}

	botToken := cfg.TelegramToken
	text := fmt.Sprintf("%s\nОписание: %s\nСтатус: %s\n\n", task.Title, task.Description, "в процессе")

	// forms url
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s",
		botToken, chatID, url.QueryEscape(text))

	// get, because request without body
	resp, err := http.Get(apiURL)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса в Telegram:", err)
		return err
	}
	defer resp.Body.Close()

	// check response
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Сообщение успешно отправлено!")
	} else {
		fmt.Printf("Не удалось отправить сообщение, статус: %s\n", resp.Status)
	}

	return nil
}
