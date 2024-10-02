package router

import (
	"net/http"
	"todo/internal/todo/middleware/middleware"
	"todo/internal/todo/transport/http/handler"

	"github.com/go-chi/chi/v5"
)

func New(todoHandler *handler.TodoHandler) http.Handler {
	r := chi.NewRouter()

	r.Route("/user", func(r chi.Router) {
		r.Post("/register", todoHandler.RegisterNewUser)
		r.Post("/login", todoHandler.AuthorizateUser)
		r.With(middleware.JWT).Get("/", todoHandler.GetAuthUser)
		r.With(middleware.JWT).Delete("/logout", todoHandler.UserLogout)
	})

	// Routes for boards
	r.Route("/boards", func(r chi.Router) {
		r.Use(middleware.JWT)
		r.Get("/", todoHandler.GetAllBoards)
		r.Get("/{id}", todoHandler.GetBoard)
		r.Post("/", todoHandler.SetBoard)
		r.Put("/{id}", todoHandler.UpdateBoard)
		r.Delete("/{id}", todoHandler.DeleteBoard)
		r.Post("/{id}", todoHandler.User2Board)
	})

	// Routes for tasks
	r.Route("/tasks", func(r chi.Router) {
		r.Use(middleware.JWT)
		r.Get("/", todoHandler.GetAllTasks)
		r.Get("/{id}", todoHandler.GetTask)
		r.Post("/", todoHandler.SetTask)
		r.Put("/{id}", todoHandler.UpdateTask)
		r.Delete("/{id}", todoHandler.DeleteTask)
	})

	r.Route("/status", func(r chi.Router) {
		r.Use(middleware.JWT)
		r.Post("/", todoHandler.SetStatus)
		r.Delete("/", todoHandler.DeleteStatus)
	})

	return r
}
