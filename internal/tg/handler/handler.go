package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"todo/internal/tg/dto"
	"todo/internal/tg/utils"

	"go.uber.org/zap"
)

type TgHandler struct {
	service TgHandlerer
	logger  *zap.Logger
}

type TgHandlerer interface {
	CreateTask(message string, chatID int64) error
	Scheduler(message string, chatID int64) error
}

func New(t TgHandlerer, logger *zap.Logger) TgHandler {
	return TgHandler{
		service: t,
		logger:  logger,
	}
}

// Handler для создания задачи
func (t *TgHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task dto.TaskDtoChatID
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	message := fmt.Sprintf("%s\nОписание: %s\nСтатус: в процессе", task.Title, task.Description)

	err := t.service.CreateTask(message, task.ChatId)
	if err != nil {
		http.Error(w, "No tg user", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Handler для обработки расписания
func (t *TgHandler) Scheduler(w http.ResponseWriter, r *http.Request) {
	var mess []dto.MessDto
	if err := json.NewDecoder(r.Body).Decode(&mess); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	chatID, message := utils.FormatTasksMessage(mess)
	if chatID == nil {
		http.Error(w, "chatID is invalid", http.StatusBadRequest)
		return
	}

	err := t.service.Scheduler(message, *chatID)
	if err != nil {
		http.Error(w, "No tg user", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
