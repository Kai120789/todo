package config

import (
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	TelegramToken string
	ToDoAppURL    string
	LogLevel      string
}

func GetConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}

	if telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN"); telegramToken != "" {
		cfg.TelegramToken = telegramToken
	} else {
		cfg.TelegramToken = zapcore.ErrorLevel.String()
	}

	if toDoAppURL := os.Getenv("TODO_APP_URL"); toDoAppURL != "" {
		cfg.ToDoAppURL = toDoAppURL
	} else {
		cfg.ToDoAppURL = "0.0.0.0:8080"
	}

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		cfg.LogLevel = envLogLevel
	} else {
		cfg.LogLevel = zapcore.ErrorLevel.String()
	}

	return cfg, nil
}
