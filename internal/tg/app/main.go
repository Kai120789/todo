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
		fmt.Println(err.Error())
	}

	// init logger
	zapLog, err := logger.New(cfg.LogLevel)
	if err != nil {
		fmt.Println(err.Error())
	}

	log := zapLog.ZapLogger

	// bot init
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal("error init bot", zap.Error(err))
	}
	log.Info(fmt.Sprintf("Authorized on account %s", bot.Self.UserName))

	go func() {
		for range time.Tick(1 * time.Second) {
			gocron.RunPending()
		}
	}()

	serv := service.New(log, bot)

	h := handler.New(serv, log)

	r := chi.NewRouter()

	srv := &http.Server{
		Addr:    cfg.TgAddress,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server", zap.Error(err))
		}
	}()

	r.Post("/create-task", h.CreateTask)
	r.Post("/scheduler", h.Scheduler)

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
				err := api.AddChatID(tgUsername, chatID, cfg.ToDoAppURL)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(chatID, "Ошибка при регистрации. Попробуйте снова."))
					return
				}

				bot.Send(tgbotapi.NewMessage(chatID, "Вы зарегистрированы!"))
			}
		}

	}

}
