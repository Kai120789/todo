package config

import (
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	TelegramToken string
	DBDSN         string
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

	if envDBConn := os.Getenv("DBDSN"); envDBConn != "" {
		cfg.DBDSN = envDBConn
	} else {
		cfg.DBDSN = zapcore.ErrorLevel.String()
	}

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		cfg.LogLevel = envLogLevel
	} else {
		cfg.LogLevel = zapcore.ErrorLevel.String()
	}

	return cfg, nil
}
