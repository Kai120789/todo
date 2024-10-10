package app

import (
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

	// get config
	cfg, err := config.GetConfig()
	if err != nil {
		zap.S().Fatalf("get config error", zap.Error(err))
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

	defer dbConn.Close()

	// create storage
	db := storage.New(dbConn, log)

	// create service
	s := services.New(services.Storager{
		BoardsStorager:   &db.BoardsStorage,
		StatusesStorager: &db.StatusesStorage,
		TasksStorager:    &db.TasksStorage,
		UserStorager:     &db.UserStorage,
	}, log)

	s.TasksService.StartScheduler()

	// init handler
	h := handler.New(handler.TodoService{
		BoardsService:   &s.BoardsService,
		StatusesService: &s.StatusesService,
		TasksService:    &s.TasksService,
		UserService:     &s.UserService,
	}, log)

	// init router
	r := router.New(&h)

	// start http-server
	log.Info("starting server", zap.String("address", cfg.ServerAddress))

	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: r,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", zap.Error(err))
	}
}
