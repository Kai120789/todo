package config

import (
	"flag"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	ServerAddress  string
	DBDSN          string
	LogLevel       string
	SecretKey      string
	TelegramToken  string
	TelegramAppURL string
}

func GetConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		cfg.ServerAddress = envRunAddr
	} else {
		flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "address and port to run server")
	}

	if envDBConn := os.Getenv("DBDSN"); envDBConn != "" {
		cfg.DBDSN = envDBConn
	} else {
		flag.StringVar(&cfg.DBDSN, "d", "", "DBDSN for database")
	}

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		cfg.LogLevel = envLogLevel
	} else {
		cfg.LogLevel = zapcore.ErrorLevel.String()
	}

	if secretKey := os.Getenv("SECRET_KEY"); secretKey != "" {
		cfg.SecretKey = secretKey
	} else {
		cfg.SecretKey = zapcore.ErrorLevel.String()
	}

	if telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN"); telegramToken != "" {
		cfg.TelegramToken = telegramToken
	} else {
		cfg.TelegramToken = zapcore.ErrorLevel.String()
	}

	if telegramAppURL := os.Getenv("TODO_APP_URL"); telegramAppURL != "" {
		cfg.TelegramAppURL = telegramAppURL
	} else {
		cfg.TelegramAppURL = zapcore.ErrorLevel.String()
	}

	flag.Parse()

	return cfg, nil
}
