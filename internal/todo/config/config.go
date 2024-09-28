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
	DbUser          string
	DbPassword      string
	DbName          string
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

	if DbUser := os.Getenv("POSTGRES_USER"); DbUser != "" {
		cfg.DbUser = DbUser
	} else {
		flag.StringVar(&cfg.DbUser, "a", "yourdbuser", "user for connect postgres")
	}

	if DbPassword := os.Getenv("POSTGRES_PASSWORD"); DbPassword != "" {
		cfg.DbPassword = DbPassword
	} else {
		flag.StringVar(&cfg.DbPassword, "a", "yourdbpassword", "password for connect postgres")
	}

	if DbName := os.Getenv("POSTGRES_DB"); DbName != "" {
		cfg.DbName = DbName
	} else {
		flag.StringVar(&cfg.DbName, "a", "dbname", "postgres dbname")
	}

	flag.Parse()

	return cfg, nil
}
