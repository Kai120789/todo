package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"todo/internal/todo/dto"
	"todo/internal/todo/models"
	"todo/internal/todo/utils/tokens"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type TodoHandler struct {
	service TodoHandlerer
	logger  *zap.Logger
}

type TodoHandlerer interface {
	SetBoard(body dto.PostBoardDto) (*models.Board, error)
	GetAllBoards() ([]models.Board, error)
	GetBoard(id uint) (*models.Board, error)
	UpdateBoard(body dto.PostBoardDto, id uint) error
	DeleteBoard(id string) error

	User2Board(body dto.PostUser2BoardDto) error

	SetTask(body dto.PostTaskDto) error
	GetTask(id uint) (*models.Task, error)
	GetAllTasks() ([]models.Task, error)
	UpdateTask(body dto.PostTaskDto, id uint) error
	DeleteTask(id string) error

	SetStatus(body dto.PostStatusDto) error
	DeleteStatus(id string) error

	RegisterNewUser(body dto.PostUserDto) (*models.UserToken, error)
	AuthorizateUser(body dto.PostUserDto) (*models.UserToken, uint, error)
	WriteRefreshToken(userId uint, refreshTokenValue string) error
	GetAuthUser(id uint) (*models.UserToken, error)
	UserLogout(id uint) error
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

	boardRet, err := h.service.SetBoard(board)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(boardRet)
}

// Get all boards
func (h *TodoHandler) GetAllBoards(w http.ResponseWriter, r *http.Request) {
	boards, err := h.service.GetAllBoards()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(board)
}

// Update a board
func (h *TodoHandler) UpdateBoard(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var board dto.PostBoardDto
	if err := json.NewDecoder(r.Body).Decode(&board); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateBoard(board, uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(board)
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

// Add user to board
func (h *TodoHandler) User2Board(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var u2b dto.PostUser2BoardDto
	if err := json.NewDecoder(r.Body).Decode(&u2b); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	u2b.BoardId = id

	if err := h.service.User2Board(u2b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(u2b)
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
	json.NewEncoder(w).Encode(task)
}

// Get all tasks
func (h *TodoHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
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

	if err := h.service.UpdateTask(task, uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
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
	json.NewEncoder(w).Encode(status)
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

func (h *TodoHandler) RegisterNewUser(w http.ResponseWriter, r *http.Request) {
	var user dto.PostUserDto
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if _, err := h.service.RegisterNewUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *TodoHandler) AuthorizateUser(w http.ResponseWriter, r *http.Request) {
	var user dto.PostUserDto
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	_, userID, err := h.service.AuthorizateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	accessTokenValue, err := tokens.GenerateJWT(userID, time.Now().Add(15*time.Minute))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	refreshTokenValue, err := tokens.GenerateJWT(userID, time.Now().Add(2*time.Hour*24*30))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	accessTokenCokie := http.Cookie{
		Name:     "access_token",
		Value:    accessTokenValue,
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		Secure:   false,
	}

	refreshTokenCokie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshTokenValue,
		Path:     "/",
		Expires:  time.Now().Add(2 * time.Hour * 24 * 30),
		HttpOnly: true,
		Secure:   false,
	}

	http.SetCookie(w, &accessTokenCokie)
	http.SetCookie(w, &refreshTokenCokie)

	err = h.service.WriteRefreshToken(userID, refreshTokenValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accessTokenValue)
}

func (h *TodoHandler) GetAuthUser(w http.ResponseWriter, r *http.Request) {
	var user dto.PostUserDto
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	_, userID, err := h.service.AuthorizateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	user2, err := h.service.GetAuthUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user2 == nil {
		http.Error(w, "No active user", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userID)
}

func (h *TodoHandler) UserLogout(w http.ResponseWriter, r *http.Request) {
	var user dto.PostUserDto
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	_, userID, err := h.service.AuthorizateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = h.service.UserLogout(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	expiredCookie := time.Now().Add(-1 * time.Hour)

	accessTokenCokie := http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		Expires:  expiredCookie,
		HttpOnly: true,
		Secure:   false, // Используйте true, если у вас HTTPS
	}

	refreshTokenCokie := http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		Expires:  expiredCookie,
		HttpOnly: true,
		Secure:   false, // Используйте true, если у вас HTTPS
	}

	http.SetCookie(w, &accessTokenCokie)
	http.SetCookie(w, &refreshTokenCokie)

	w.WriteHeader(http.StatusNoContent)
}
