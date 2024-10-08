package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
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
	SendDailyReports(tasks []models.Task, chatID int64, status int) error
}

func New(t TgHandlerer, logger *zap.Logger) TgHandler {
	return TgHandler{
		service: t,
		logger:  logger,
	}
}

// Handler для создания задачи
func (t *TgHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	// Валидация и парсинг входящих данных
	var requestData struct {
		Task   models.Task `json:"task"`
		ChatID int64       `json:"chatID"`
	}

	// Декодируем тело запроса в структуру requestData
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		t.logger.Error("Invalid request data", zap.Error(err))
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	// Валидация chatID
	if requestData.ChatID == 0 {
		t.logger.Error("Chat ID is missing or invalid")
		http.Error(w, "Chat ID is required", http.StatusBadRequest)
		return
	}

	// Вызов соответствующей сервисной функции
	err = t.service.CreateTask(&requestData.Task, requestData.ChatID)
	if err != nil {
		t.logger.Error("Failed to create task", zap.Error(err))
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	// Успешный ответ
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Task created successfully"))
}

// Handler для обработки расписания
func (t *TgHandler) Scheduler(w http.ResponseWriter, r *http.Request) {
	// Валидация и парсинг входящих данных (массив задач)
	var tasks []models.Task
	err := json.NewDecoder(r.Body).Decode(&tasks)
	if err != nil {
		t.logger.Error("Invalid task data", zap.Error(err))
		http.Error(w, "Invalid task data", http.StatusBadRequest)
		return
	}

	chatIDStr := r.FormValue("chatID")
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		t.logger.Error("Invalid chat ID", zap.Error(err))
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	// Вызов соответствующей сервисной функции
	err = t.service.Scheduler(tasks, chatID)
	if err != nil {
		t.logger.Error("Failed to send scheduler", zap.Error(err))
		http.Error(w, "Failed to send scheduler", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Scheduler sent successfully"))
}
