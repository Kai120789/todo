package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"todo/internal/repositories"

	"github.com/gorilla/mux"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:  "success",
		Message: "Сервер работает",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

var taskRepo = repositories.NewTaskRepository()

// GetTasks — получение всех задач для текущего пользователя
func GetTasks(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int) // Извлекаем userID из контекста (добавляется в middleware)

	tasks, err := taskRepo.GetTasksByUserID(userID)
	if err != nil {
		log.Printf("Error fetching tasks: %v", err)
		http.Error(w, "Unable to fetch tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// GetTaskByID — получение задачи по ID
func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, err := taskRepo.GetTaskByID(userID, taskID)
	if err != nil {
		log.Printf("Error fetching task: %v", err)
		http.Error(w, "Unable to fetch task", http.StatusInternalServerError)
		return
	}

	if task == nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}
