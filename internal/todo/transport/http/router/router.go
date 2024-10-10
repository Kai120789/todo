package router

import (
	"net/http"
	"todo/internal/todo/middleware/middleware"
	"todo/internal/todo/transport/http/handler"

	"github.com/go-chi/chi/v5"
)

func New(todoHandler *handler.TodoHandler) http.Handler {
	r := chi.NewRouter()

	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", todoHandler.RegisterNewUser)                 // register new user
		r.Post("/login", todoHandler.AuthorizateUser)                    // login user
		r.With(middleware.JWT).Get("/", todoHandler.GetAuthUser)         // get active user, need jwt
		r.With(middleware.JWT).Delete("/logout", todoHandler.UserLogout) // logout user, need jwt
	})

	// Routes for boards
	r.Route("/api/boards", func(r chi.Router) {
		r.Use(middleware.JWT)                      // need jwt for all methods
		r.Get("/", todoHandler.GetAllBoards)       // get all boards
		r.Get("/{id}", todoHandler.GetBoard)       // get board with id
		r.Post("/", todoHandler.SetBoard)          // add new board
		r.Put("/{id}", todoHandler.UpdateBoard)    // update board
		r.Delete("/{id}", todoHandler.DeleteBoard) // delete board
		r.Post("/{id}", todoHandler.User2Board)    // add user to board
	})

	// Routes for tasks
	r.Route("/api/tasks", func(r chi.Router) {
		r.Use(middleware.JWT)                     // need jwt for all methods
		r.Get("/", todoHandler.GetAllTasks)       // get all tasks
		r.Get("/{id}", todoHandler.GetTask)       // get task with id
		r.Post("/", todoHandler.SetTask)          // add new task
		r.Put("/{id}", todoHandler.UpdateTask)    // update task
		r.Delete("/{id}", todoHandler.DeleteTask) // delete task
	})

	// Routes for statuses
	r.Route("/api/status", func(r chi.Router) {
		r.Use(middleware.JWT)                   // need jwt for all methods
		r.Post("/", todoHandler.SetStatus)      // add new status
		r.Delete("/", todoHandler.DeleteStatus) // delete status
	})

	r.Post("/add-chat-id", todoHandler.AddChatID) // add chatID to table users

	return r
}
