package service

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

type TgService struct {
	logger *zap.Logger
	bot    *tgbotapi.BotAPI
}

type TgServiceer interface {
	CreateTask(message string, chatID int64) error
	Scheduler(message string, chatID int64) error
}

// Конструктор для TgService
func New(logger *zap.Logger, bot *tgbotapi.BotAPI) *TgService {
	return &TgService{
		logger: logger,
		bot:    bot,
	}
}

// Создание задачи и отправка сообщения в Telegram
func (s *TgService) CreateTask(message string, chatID int64) error {

	s.bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Добавлена новая задача:\n\n%s", message)))

	return nil
}

// Отправка расписания в Telegram
func (s *TgService) Scheduler(message string, chatID int64) error {

	s.bot.Send(tgbotapi.NewMessage(chatID, message))

	return nil
}
