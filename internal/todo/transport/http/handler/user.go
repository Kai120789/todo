package handler

import (
	"encoding/json"
	"net/http"
	"time"
	"todo/internal/todo/dto"
	"todo/internal/todo/models"
	"todo/internal/todo/utils/tokens"

	"go.uber.org/zap"
)

type UserHandler struct {
	service UserHandlerer
	logger  *zap.Logger
}

type UserHandlerer interface {
	RegisterNewUser(body dto.PostUserDto) (*models.UserToken, error)
	AuthorizateUser(body dto.PostUserDto) (*uint, error)
	WriteRefreshToken(userId uint, refreshTokenValue string) error
	GetAuthUser(id uint) (*models.UserToken, error)
	UserLogout(id uint) error
	AddChatID(tgName string, chatID int64) error
}

func NewUserHandler(t UserHandlerer, logger *zap.Logger) UserHandler {
	return UserHandler{
		service: t,
		logger:  logger,
	}
}

// Register new user
func (h *UserHandler) RegisterNewUser(w http.ResponseWriter, r *http.Request) {
	var user dto.PostUserDto
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if user.Username == "" {
		http.Error(w, "username cannot be empty", http.StatusBadRequest)
		return
	}

	if _, err := h.service.RegisterNewUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Login user
func (h *UserHandler) AuthorizateUser(w http.ResponseWriter, r *http.Request) {
	var user dto.PostUserDto
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if user.Username == "" {
		http.Error(w, "username cannot be empty", http.StatusBadRequest)
		return
	}

	userID, err := h.service.AuthorizateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	accessTokenValue, err := tokens.GenerateJWT(*userID, time.Now().Add(15*time.Minute))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	refreshTokenValue, err := tokens.GenerateJWT(*userID, time.Now().Add(2*time.Hour*24*30))
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

	err = h.service.WriteRefreshToken(*userID, refreshTokenValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accessTokenValue)
}

// Get active user
func (h *UserHandler) GetAuthUser(w http.ResponseWriter, r *http.Request) {
	var user dto.PostUserDto
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if user.Username == "" {
		http.Error(w, "username cannot be empty", http.StatusBadRequest)
		return
	}

	userID, err := h.service.AuthorizateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	user2, err := h.service.GetAuthUser(*userID)
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

// Logout user
func (h *UserHandler) UserLogout(w http.ResponseWriter, r *http.Request) {
	var user dto.PostUserDto
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if user.Username == "" {
		http.Error(w, "username cannot be empty", http.StatusBadRequest)
		return
	}

	userID, err := h.service.AuthorizateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = h.service.UserLogout(*userID)
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
		Secure:   false,
	}

	refreshTokenCokie := http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		Expires:  expiredCookie,
		HttpOnly: true,
		Secure:   false,
	}

	http.SetCookie(w, &accessTokenCokie)
	http.SetCookie(w, &refreshTokenCokie)

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) AddChatID(w http.ResponseWriter, r *http.Request) {
	var user dto.PostUserDto
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err := h.service.AddChatID(user.TgName, user.ChatID)
	if err != nil {
		http.Error(w, "No tg user", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
