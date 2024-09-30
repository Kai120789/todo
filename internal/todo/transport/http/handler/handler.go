package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"todo/internal/todo/dto"
	"todo/internal/todo/models"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type TodoHandler struct {
	service TodoHandlerer
	logger  *zap.Logger
}

type TodoHandlerer interface {
	SetBoard(body dto.PostBoardDto) error
	GetAllBoards() ([]models.Board, error)
	GetBoard(id uint) (*models.Board, error)
	UpdateBoard(body dto.PostBoardDto) error
	DeleteBoard(id string) error
	SetTask(body dto.PostTaskDto) error
	GetTask(id uint) (*models.Task, error)
	GetAllTasks() ([]models.Task, error)
	UpdateTask(body dto.PostTaskDto) error
	DeleteTask(id string) error
	SetStatus(body dto.PostStatusDto) error
	DeleteStatus(id string) error
}

func New(t TodoHandlerer, logger *zap.Logger) TodoHandler {
	return TodoHandler{
		service: t,
		logger:  logger,
	}
}

func (h *TodoHandler) SetBoard(w http.ResponseWriter, r *http.Request) {
	var board dto.PostBoardDto
	if err := json.NewDecoder(r.Body).Decode(&board); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := h.service.SetBoard(board); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Get all boards
func (h *TodoHandler) GetAllBoards(w http.ResponseWriter, r *http.Request) {
	boards, err := h.service.GetAllBoards()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(boards)
}

// Get a specific board
func (h *TodoHandler) GetBoard(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid ID: %s, error: %v", idStr, err), http.StatusBadRequest)
		return
	}

	board, err := h.service.GetBoard(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if board == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(board)
}

// Update a board
func (h *TodoHandler) UpdateBoard(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var board dto.PostBoardDto
	if err := json.NewDecoder(r.Body).Decode(&board); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	board.ID = id

	if err := h.service.UpdateBoard(board); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Delete a board
func (h *TodoHandler) DeleteBoard(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.DeleteBoard(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Create a new task
func (h *TodoHandler) SetTask(w http.ResponseWriter, r *http.Request) {
	var task dto.PostTaskDto
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := h.service.SetTask(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Get all tasks
func (h *TodoHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.GetAllTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// Get a specific task
func (h *TodoHandler) GetTask(w http.ResponseWriter, r *http.Request) {
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
func (h *TodoHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var task dto.PostTaskDto
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	task.ID = id

	if err := h.service.UpdateTask(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Delete a task
func (h *TodoHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.DeleteTask(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SetStatus - создаёт новый статус
func (h *TodoHandler) SetStatus(w http.ResponseWriter, r *http.Request) {
	var status dto.PostStatusDto
	if err := json.NewDecoder(r.Body).Decode(&status); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := h.service.SetStatus(status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// DeleteStatus - удаляет существующий статус
func (h *TodoHandler) DeleteStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.DeleteStatus(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
