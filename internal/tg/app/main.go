package app

import (
	"fmt"
	"time"
	"todo/internal/tg/config"
	tgservice "todo/internal/tg/service"
	tgstorage "todo/internal/tg/storage"
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
	dbConn, err := storage.Connection("postgres://postgres:123456@localhost:5431/taskdb?sslmode=disable")
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

	store := tgstorage.New(dbConn, log)

	serv := tgservice.New(store, log)

	// all tasks at 00:00
	gocron.Every(1).Day().At("00:00").Do(func() {
		users, err := store.GetAllUsers()
		if err != nil {
			log.Error("Ошибка при получении пользователей", zap.Error(err))
			return
		}

		for _, user := range users {
			message, _, err := serv.GetMyTasks(user.TgName)
			if err != nil {
				log.Error("Ошибка получения задач для пользователя", zap.String("tgName", user.TgName), zap.Error(err))
				continue
			}
			bot.Send(tgbotapi.NewMessage(user.ChatID, message))

			message, _, err = serv.GetMyEndedTasks(user.TgName)
			if err != nil {
				log.Error("Ошибка получения задач для пользователя", zap.String("tgName", user.TgName), zap.Error(err))
				continue
			}
			bot.Send(tgbotapi.NewMessage(user.ChatID, message))

			err = serv.ChangeEndedTasksStatus()
			if err != nil {
				log.Error("Ошибка побновления статуса", zap.String("tgName", user.TgName), zap.Error(err))
				return
			}
		}
	})

	go func() {
		for range time.Tick(1 * time.Second) {
			gocron.RunPending()
		}
	}()

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
			chatID := update.Message.Chat.ID
			tgUsername := update.Message.From.UserName

			if tgUsername == "" {
				bot.Send(tgbotapi.NewMessage(chatID, "У вас не установлен Telegram username"))
				return
			}

			switch update.Message.Command() {
			case "start":
				err := store.RegisterUser(update.Message.From.ID, tgUsername, chatID)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(chatID, "Ошибка при регистрации"))
					return
				}
				bot.Send(tgbotapi.NewMessage(chatID, "Вы зарегистрированы!"))

			case "my_tasks":
				message, _, err := serv.GetMyTasks(tgUsername)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(chatID, "Ошибка поиска задач"))
					return
				}
				bot.Send(tgbotapi.NewMessage(chatID, message))
			}
		}

	}

}
