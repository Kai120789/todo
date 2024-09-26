package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"todo/internal/config"
	"todo/internal/db"
	"todo/internal/lib/logger/sl"
	"todo/internal/routers"

	"github.com/joho/godotenv"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	cfg := config.MustLoad()

	log := settupLogger(cfg.Env)

	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	fmt.Println(cfg)

	err := db.InitDB(cfg.User, cfg.Password, cfg.Name, cfg.Host, int(cfg.Port))
	if err != nil {
		log.Error("Ошибка инициализации базы данных", sl.Err(err))
	}

	defer db.CloseDB()

	router := routers.SetupRouter()

	log.Info("Запуск сервера на http://localhost:8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Error("Ошибка запуска сервера", sl.Err(err))
	}

}

func settupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log

}
