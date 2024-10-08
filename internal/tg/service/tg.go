package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"todo/internal/tg/config"
	"todo/internal/tg/models"

	"go.uber.org/zap"
)

type TgService struct {
	logger *zap.Logger
}

type TgServiceer interface {
	CreateTask(task *models.Task, chatID int64) error
	Scheduler(tasks []models.Task, chatID int64) error
	SendDailyReports(tasks []models.Task, chatID int64, status int) error
}

// Конструктор для TgService
func New(logger *zap.Logger) *TgService {
	return &TgService{
		logger: logger,
	}
}

// Создание задачи и отправка сообщения в Telegram
func (s *TgService) CreateTask(task *models.Task, chatID int64) error {
	cfg, err := config.GetConfig()
	if err != nil {
		s.logger.Fatal("Error loading config", zap.Error(err))
	}

	// Формируем текст сообщения
	message := fmt.Sprintf("New Task: %s\nDetails: %s\nStatus: %d", task.Title, task.Description, task.StatusId)

	// Формируем URL для отправки сообщения через Telegram API
	botUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s", cfg.TelegramToken, chatID, message)

	// Формируем данные для отправки в JSON формате
	taskData := map[string]interface{}{
		"chat_id": chatID,
		"text":    message,
	}

	jsonData, err := json.Marshal(taskData)
	if err != nil {
		s.logger.Error("Failed to marshal task data", zap.Error(err))
		return err
	}

	// Выполняем POST запрос
	resp, err := http.Post(botUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		s.logger.Error("Failed to send task to Telegram", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	// Проверка статуса ответа
	if resp.StatusCode == http.StatusOK {
		s.logger.Info("Task successfully sent to Telegram")
	} else {
		zap.S().Error("Failed to send task, status: %s", resp.Status)
	}

	return nil
}

// Отправка расписания в Telegram
func (s *TgService) Scheduler(tasks []models.Task, chatID int64) error {
	cfg, err := config.GetConfig()
	if err != nil {
		s.logger.Fatal("Error loading config", zap.Error(err))
	}

	// Формируем сообщение для отправки
	message := formatSchedulerMessage(tasks)

	// Формируем URL для отправки сообщения через Telegram API
	botUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", cfg.TelegramToken)

	// Параметры запроса
	params := url.Values{}
	params.Add("chat_id", fmt.Sprintf("%d", chatID))
	params.Add("text", message)

	// Выполняем запрос на отправку сообщения
	resp, err := http.PostForm(botUrl, params)
	if err != nil {
		s.logger.Error("Failed to send scheduler to Telegram", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	// Проверка статуса ответа
	if resp.StatusCode == http.StatusOK {
		s.logger.Info("Scheduler successfully sent to Telegram")
	} else {
		zap.S().Error("Failed to send scheduler, status: %s", resp.Status)
	}

	return nil
}

// Отправка ежедневных отчетов в Telegram
func (s *TgService) SendDailyReports(tasks []models.Task, chatID int64, status int) error {
	cfg, err := config.GetConfig()
	if err != nil {
		s.logger.Fatal("Error loading config", zap.Error(err))
	}

	// Формирование сообщения в зависимости от статуса
	var message string
	if status == 1 {
		message = formatTasksMessage(tasks)
	} else {
		message = formatEndedTasksMessage(tasks)
	}

	// Формируем URL для отправки сообщения через Telegram API
	botUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", cfg.TelegramToken)

	// Параметры запроса
	params := url.Values{}
	params.Add("chat_id", fmt.Sprintf("%d", chatID))
	params.Add("text", message)

	// Выполняем запрос на отправку сообщения
	resp, err := http.PostForm(botUrl, params)
	if err != nil {
		s.logger.Error("Failed to send daily report to Telegram", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	// Проверка статуса ответа
	if resp.StatusCode == http.StatusOK {
		s.logger.Info("Daily report successfully sent to Telegram")
	} else {
		zap.S().Error("Failed to send daily report, status: %s", resp.Status)
	}

	return nil
}

// Вспомогательные функции для форматирования сообщений
func formatSchedulerMessage(tasks []models.Task) string {
	var message string
	for _, task := range tasks {
		message += fmt.Sprintf("Task: %s, Status: %d\n", task.Title, task.StatusId)
	}
	return message
}

func formatEndedTasksMessage(tasks []models.Task) string {
	var message string
	for _, task := range tasks {
		message += fmt.Sprintf("Completed Task: %s\n", task.Title)
	}
	return message
}

func formatTasksMessage(tasks []models.Task) string {
	var message string
	for _, task := range tasks {
		message += fmt.Sprintf("Task: %s\n", task.Title)
	}
	return message
}
