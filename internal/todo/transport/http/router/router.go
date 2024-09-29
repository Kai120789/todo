package router

import (
	"net/http"
	"todo/internal/todo/transport/http/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(todoHandler *handler.TodoHandler) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	/*r.Route("/user", func(r chi.Router) {
		r.Post("/register", todoHandler.RegisterNewUser)
		r.Post("/login", todoHandler.AuthorizateUser)
		r.Get("/", todoHandler.GetAuthUser)
	})*/

	// Routes for boards
	r.Route("/boards", func(r chi.Router) {
		r.Get("/", todoHandler.GetAllBoards)
		r.Get("/{id}", todoHandler.GetBoard)
		r.Post("/", todoHandler.SetBoard)
		r.Put("/{id}", todoHandler.UpdateBoard)
		r.Delete("/{id}", todoHandler.DeleteBoard)
	})

	// Routes for tasks
	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", todoHandler.GetAllTasks)
		r.Get("/tasks/{id}", todoHandler.GetTask)
		r.Post("/", todoHandler.GetAllTasks)
		r.Put("/tasks/{id}", todoHandler.UpdateTask)
		r.Delete("/tasks/{id}", todoHandler.DeleteTask)
	})

	r.Post("/status", todoHandler.SetStatus)
	r.Delete("/status", todoHandler.DeleteStatus)

	return r
}
