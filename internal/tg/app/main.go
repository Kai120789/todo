package app

import (
	"fmt"
	"time"
	"todo/internal/tg/config"
	"todo/internal/todo/storage"
	"todo/pkg/logger"

	"github.com/jasonlvhit/gocron"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

func StartTgBot() {

	// init config
	cfg, err := config.GetConfig()
	if err != nil {
		zap.S().Fatal("error load config", zap.Error(err))
	}

	// init logger
	zapLog, err := logger.New(cfg.LogLevel)
	if err != nil {
		zap.S().Fatalf("init logger error", zap.Error(err))
	}

	log := zapLog.ZapLogger

	// connect to postgres db
	dbConn, err := storage.Connection(cfg.DBDSN)
	if err != nil {
		log.Fatal("error connect to db", zap.Error(err))
	}

	// bot init
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal("error init bot", zap.Error(err))
	}
	log.Info(fmt.Sprintf("Authorized on account %s", bot.Self.UserName))

	_ = log

	defer dbConn.Close()

	// Запуск задачи в 00:00
	//gocron.Every(1).Day().At("00:00").Do(sendDailyReport)

	// Запуск планировщика
	go func() {
		for range time.Tick(1 * time.Second) {
			gocron.RunPending()
		}
	}()

	// Обработка обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return
	}

	for update := range updates {
		if update.Message == nil { // ignore non-Message Updates
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				// Логика для регистрации пользователя
				//registerUser(update.Message.From.ID, update.Message.From.UserName)
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Вы зарегистрированы!"))
			case "add_task":
				// Пример добавления задачи
				//addTask(update.Message.From.ID, "Пример задачи")
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Задача добавлена!"))
			case "add_board":
				// Пример добавления доски
				//addBoard(update.Message.From.ID, "Пример доски")
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Доска добавлена!"))
			}
		}
	}
}
