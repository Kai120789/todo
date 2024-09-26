package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"todo/internal/models"
	"todo/internal/repositories"
	"todo/internal/utils"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}
	user.PasswordHash = hashedPassword

	id, err := repositories.CreateUser(user)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	user.ID = id
	json.NewEncoder(w).Encode(user)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var credentials models.User
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	log.Printf("Attempting login for user: %s", credentials.Username)

	user, err := repositories.GetUserByUsername(credentials.Username)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	log.Printf("User found: %+v", user)

	if !utils.CheckPasswordHash(credentials.Password, user.PasswordHash) {
		log.Printf("Invalid password for user: %s", credentials.Username)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	log.Printf("User %s logged in successfully", credentials.Username)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
