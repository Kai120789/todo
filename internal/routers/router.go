package routers

import (
	"todo/internal/handlers"
	"todo/internal/middleware"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/health", handlers.HealthCheck).Methods("GET")
	router.HandleFunc("/user/register", handlers.Register).Methods("POST")
	router.HandleFunc("/user/login", handlers.Login).Methods("POST")

	secured := router.PathPrefix("/api").Subrouter()
	secured.Use(middleware.AuthMiddleware)
	secured.HandleFunc("/tasks", handlers.GetTasks).Methods("GET")

	secured.HandleFunc("/boards", handlers.CreateBoard).Methods("POST")
	secured.HandleFunc("/boards", handlers.GetBoards).Methods("GET")
	secured.HandleFunc("/boards/{id}", handlers.GetBoardByID).Methods("GET")
	secured.HandleFunc("/boards/{id}", handlers.UpdateBoard).Methods("PUT")
	secured.HandleFunc("/boards/{id}", handlers.DeleteBoard).Methods("DELETE")

	return router
}
