package app

import (
	"fmt"
	"todo/internal/todo/config"
	"todo/internal/todo/storage"
	"todo/pkg/logger"

	"go.uber.org/zap"
)

func StartServer() {

	// Получаем конфиг
	cfg, err := config.GetConfig()
	if err != nil {
		zap.S().Fatalf("get config error", zap.Error(err))
	}

	fmt.Println(cfg)

	// инициализируем логгер
	log, err := logger.New(cfg.LogLevel)
	if err != nil {
		zap.S().Fatalf("init logger error", zap.Error(err))
	}

	fmt.Println(log)

	// подключаемся к бд
	dbConn, err := storage.Connection(cfg.DSN, log.ZapLogger)
	if err != nil {
		log.ZapLogger.Fatal("error connect to db", zap.Error(err))
	}

	_ = dbConn

	defer dbConn.Close()

}
