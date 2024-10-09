package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"todo/internal/tg/dto"

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
	fmt.Println(1)

	var task dto.TaskDtoChatID
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	message := fmt.Sprintf("%s\n%s\n%d", task.Title, task.Description, task.StatusId)

	err := t.service.CreateTask(message, task.ChatId)
	if err != nil {
		http.Error(w, "No tg user", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Handler для обработки расписания
func (t *TgHandler) Scheduler(w http.ResponseWriter, r *http.Request) {
	var mess dto.Dto
	if err := json.NewDecoder(r.Body).Decode(&mess); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	message := fmt.Sprintf("%s\n\n%s", mess.Message, mess.MessageEnded)

	err := t.service.Scheduler(message, mess.ChatID)
	if err != nil {
		http.Error(w, "No tg user", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
