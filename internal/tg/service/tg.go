package service

import (
	"todo/internal/tg/models"

	"go.uber.org/zap"
)

type TgService struct {
	logger *zap.Logger
}

type TgServiceer interface {
	CreateTask(task *models.Task, chatID int64) error
	Scheduler(tasks []models.Task, chatID int64) error
}

// Конструктор для TgService
func New(logger *zap.Logger) *TgService {
	return &TgService{
		logger: logger,
	}
}

// Создание задачи и отправка сообщения в Telegram
func (s *TgService) CreateTask(task *models.Task, chatID int64) error {

}

// Отправка расписания в Telegram
func (s *TgService) Scheduler(tasks []models.Task, chatID int64) error {

}
