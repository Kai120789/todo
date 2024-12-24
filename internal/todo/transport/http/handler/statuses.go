package handler

import (
	"encoding/json"
	"net/http"
	"todo/internal/todo/dto"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type StatusesHandler struct {
	service StatusesHandlerer
	logger  *zap.Logger
}

type StatusesHandlerer interface {
	SetStatus(body dto.PostStatusDto) error
	DeleteStatus(id string) error
}

func NewStatusesHandler(t StatusesHandlerer, logger *zap.Logger) StatusesHandler {
	return StatusesHandler{
		service: t,
		logger:  logger,
	}
}

// SetStatus
func (h *StatusesHandler) SetStatus(w http.ResponseWriter, r *http.Request) {
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

// DeleteStatus
func (h *StatusesHandler) DeleteStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.DeleteStatus(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
