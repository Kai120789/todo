package app

import (
	"fmt"
	"net/http"
	"todo/internal/todo/config"
	"todo/internal/todo/services"
	"todo/internal/todo/storage"
	"todo/internal/todo/transport/http/handler"
	"todo/internal/todo/transport/http/router"
	"todo/pkg/logger"

	"go.uber.org/zap"
)

func StartServer() {

	// Получаем конфиг
	cfg, err := config.GetConfig()
	if err != nil {
		zap.S().Fatalf("get config error", zap.Error(err))
	}

	// инициализируем логгер
	zapLog, err := logger.New(cfg.LogLevel)
	if err != nil {
		zap.S().Fatalf("init logger error", zap.Error(err))
	}

	log := zapLog.ZapLogger

	// подключаемся к бд
	dbConn, err := storage.Connection(cfg.DSN)
	if err != nil {
		log.Fatal("error connect to db", zap.Error(err))
	}

	defer dbConn.Close()

	// создание хранилища
	store := storage.New(dbConn, log)

	// создание сервисного слоя
	serv := services.New(store, log)

	// инициализация хэндлера
	handl := handler.New(serv, log)

	// инициализация роутера
	r := router.New(&handl)

	_ = r

	fmt.Println(cfg)

	// настройка и запуск http-сервиса
	zap.S().Info("starting server", cfg.ServerAddress)

	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: r,
	}

	if err := srv.ListenAndServe(); err != nil {
		zap.S().Error("failed to start server")
	}

	zap.S().Error("server stopped")
}