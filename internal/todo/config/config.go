package config

import (
	"flag"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	ServerAddress        string
	FileStoragePath      string
	DBDSN                string
	envSalt              string
	LogLevel             string
	SecretKey            string
	AccessTokenTimeLife  time.Duration
	RefreshTokenTimeLife time.Duration
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

	if envSalt := os.Getenv("SALT"); envSalt != "" {
		cfg.envSalt = envSalt
	} else {
		cfg.envSalt = zapcore.ErrorLevel.String()
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

	flag.Parse()

	return cfg, nil
}
