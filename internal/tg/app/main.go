package app

import (
	"fmt"
	"net/http"
	"time"
	"todo/internal/tg/api"
	"todo/internal/tg/config"
	"todo/internal/tg/handler"
	"todo/internal/tg/service"
	"todo/pkg/logger"

	"github.com/go-chi/chi/v5"
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

	// bot init
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal("error init bot", zap.Error(err))
	}
	log.Info(fmt.Sprintf("Authorized on account %s", bot.Self.UserName))

	_ = log

	go func() {
		for range time.Tick(1 * time.Second) {
			gocron.RunPending()
		}
	}()

	serv := service.New(log)

	h := handler.New(serv, log)

	r := chi.NewRouter()
	r.Post("/create-task", h.CreateTask) //- достаем из боди дто созданной задачи, форматируем в строку и отправляем в телеграм
	r.Post("/scheduler", h.Scheduler)    //- достаем из боди массив дто выполненных задачи, форматируем в массив строк и отправляем в телеграм

	go func() {
		if err := http.ListenAndServe(":8081", r); err != nil {
			log.Fatal("Failed to start HTTP server", zap.Error(err))
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
				if api.AddChatID(tgUsername, chatID, cfg.ToDoAppURL) {
					bot.Send(tgbotapi.NewMessage(chatID, "Вы зарегистрированы!"))
				} else {
					bot.Send(tgbotapi.NewMessage(chatID, "Ошибка при регистрации. Попробуйте снова."))
				}
			}
		}

	}

}
