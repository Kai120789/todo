package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"todo/internal/todo/dto"
	"todo/internal/todo/models"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type TasksHandler struct {
	service TasksHandlerer
	logger  *zap.Logger
}

type TasksHandlerer interface {
	SetTask(body dto.PostTaskDto) error
	GetTask(id uint) (*models.Task, error)
	GetAllTasks() ([]models.Task, error)
	UpdateTask(body dto.PostTaskDto, id uint) error
	DeleteTask(id string) error
}

func NewTasksHandler(t TasksHandlerer, logger *zap.Logger) TasksHandler {
	return TasksHandler{
		service: t,
		logger:  logger,
	}
}

// Create a new task
func (h *TasksHandler) SetTask(w http.ResponseWriter, r *http.Request) {
	var task dto.PostTaskDto
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, "task title cannot be empty", http.StatusBadRequest)
		return
	}

	if err := h.service.SetTask(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// Get all tasks
func (h *TasksHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.GetAllTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

// Get a specific task
func (h *TasksHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	task, err := h.service.GetTask(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if task == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

// Update a task
func (h *TasksHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var task dto.PostTaskDto
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, "task title cannot be empty", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateTask(task, uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

// Delete a task
func (h *TasksHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.DeleteTask(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
