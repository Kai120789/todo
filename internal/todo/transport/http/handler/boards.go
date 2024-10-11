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

type BoardsHandler struct {
	service BoardsHandlerer
	logger  *zap.Logger
}

type BoardsHandlerer interface {
	SetBoard(body dto.PostBoardDto) (*models.Board, error)
	GetAllBoards() ([]models.Board, error)
	GetBoard(id uint) (*models.Board, error)
	UpdateBoard(body dto.PostBoardDto, id uint) error
	DeleteBoard(id string) error

	User2Board(body dto.PostUser2BoardDto) error
}

func NewBoardsHandler(t BoardsHandlerer, logger *zap.Logger) BoardsHandler {
	return BoardsHandler{
		service: t,
		logger:  logger,
	}
}

func (h *BoardsHandler) SetBoard(w http.ResponseWriter, r *http.Request) {
	var board dto.PostBoardDto
	if err := json.NewDecoder(r.Body).Decode(&board); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if board.Name == "" {
		http.Error(w, "Board name is empty", http.StatusBadRequest)
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
func (h *BoardsHandler) GetAllBoards(w http.ResponseWriter, r *http.Request) {
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
func (h *BoardsHandler) GetBoard(w http.ResponseWriter, r *http.Request) {
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
func (h *BoardsHandler) UpdateBoard(w http.ResponseWriter, r *http.Request) {
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

	if board.Name == "" {
		http.Error(w, "Board name is empty", http.StatusBadRequest)
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
func (h *BoardsHandler) DeleteBoard(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.DeleteBoard(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Add user to board
func (h *BoardsHandler) User2Board(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var u2b dto.PostUser2BoardDto
	if err := json.NewDecoder(r.Body).Decode(&u2b); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	u2b.BoardId = id

	if u2b.UserId == "" {
		http.Error(w, "user ID is required", http.StatusBadRequest)
		return
	}

	if err := h.service.User2Board(u2b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(u2b)
}
