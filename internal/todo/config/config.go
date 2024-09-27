package config

import (
	"flag"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	ServerAddress   string
	FileStoragePath string
	DSN             string
	envSold         string
	LogLevel        string
}

var AppConfig Config

func GetConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		cfg.ServerAddress = envRunAddr
	} else {
		flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "address and port to run server")
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		cfg.FileStoragePath = envFileStoragePath
	} else {
		flag.StringVar(&cfg.FileStoragePath, "f", "", "path to storage file")
	}

	if envDBConn := os.Getenv("DSN"); envDBConn != "" {
		cfg.DSN = envDBConn
	} else {
		flag.StringVar(&cfg.DSN, "d", "", "dsn for database")
	}

	if envSold := os.Getenv("SOLD"); envSold != "" {
		cfg.envSold = envSold
	} else {
		cfg.envSold = zapcore.ErrorLevel.String()
	}

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		cfg.LogLevel = envLogLevel
	} else {
		cfg.LogLevel = zapcore.ErrorLevel.String()
	}

	flag.Parse()

	return cfg, nil
}
