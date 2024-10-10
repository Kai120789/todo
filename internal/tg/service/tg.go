package service

import (
	"fmt"
	"todo/internal/tg/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

type TgService struct {
	logger *zap.Logger
}

type TgServiceer interface {
	CreateTask(message string, chatID int64) error
	Scheduler(message string, chatID int64) error
}

// Конструктор для TgService
func New(logger *zap.Logger) *TgService {
	return &TgService{
		logger: logger,
	}
}

// Создание задачи и отправка сообщения в Telegram
func (s *TgService) CreateTask(message string, chatID int64) error {
	cfg, err := config.GetConfig()
	if err != nil {
		zap.S().Fatal("error load config", zap.Error(err))
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		zap.S().Fatal("error init bot", zap.Error(err))
	}

	bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Добавлена новая задача:\n\n%s", message)))

	return nil
}

// Отправка расписания в Telegram
func (s *TgService) Scheduler(message string, chatID int64) error {
	cfg, err := config.GetConfig()
	if err != nil {
		zap.S().Fatal("error load config", zap.Error(err))
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		zap.S().Fatal("error init bot", zap.Error(err))
	}

	bot.Send(tgbotapi.NewMessage(chatID, message))

	return nil
}
