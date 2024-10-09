package handler

import (
	"net/http"
	"todo/internal/tg/models"

	"go.uber.org/zap"
)

type TgHandler struct {
	service TgHandlerer
	logger  *zap.Logger
}

type TgHandlerer interface {
	CreateTask(task *models.Task, chatID int64) error
	Scheduler(tasks []models.Task, chatID int64) error
}

func New(t TgHandlerer, logger *zap.Logger) TgHandler {
	return TgHandler{
		service: t,
		logger:  logger,
	}
}

// Handler для создания задачи
func (t *TgHandler) CreateTask(w http.ResponseWriter, r *http.Request) {

}

// Handler для обработки расписания
func (t *TgHandler) Scheduler(w http.ResponseWriter, r *http.Request) {

}
