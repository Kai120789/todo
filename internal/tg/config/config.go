package config

import (
	"flag"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	TelegramToken string
	ToDoAppURL    string
	LogLevel      string
	TgAddress     string
}

func GetConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}

	if telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN"); telegramToken != "" {
		cfg.TelegramToken = telegramToken
	} else {
		cfg.TelegramToken = zapcore.ErrorLevel.String()
	}

	if envRunAddr := os.Getenv("TG_ADDRESS"); envRunAddr != "" {
		cfg.TgAddress = envRunAddr
	} else {
		flag.StringVar(&cfg.TgAddress, "a", "localhost:8080", "address and port to run server")
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
