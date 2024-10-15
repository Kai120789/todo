package router

import (
	"net/http"
	"todo/internal/todo/middleware"

	"github.com/go-chi/chi/v5"
)

type TasksRouter struct{}

type TasksHandler interface {
	SetTask(w http.ResponseWriter, r *http.Request)
	GetAllTasks(w http.ResponseWriter, r *http.Request)
	GetTask(w http.ResponseWriter, r *http.Request)
	UpdateTask(w http.ResponseWriter, r *http.Request)
	DeleteTask(w http.ResponseWriter, r *http.Request)
	SendAllTasks(w http.ResponseWriter, r *http.Request)
}

func NewTasksRouter() *TasksRouter {
	return &TasksRouter{}
}

func (b *TasksRouter) TasksRoutes(r chi.Router, h TasksHandler) {
	// Routes for tasks
	r.Route("/api/tasks", func(r chi.Router) {
		r.Use(middleware.JWT)           // need jwt for all methods
		r.Get("/", h.GetAllTasks)       // get all tasks
		r.Get("/{id}", h.GetTask)       // get task with id
		r.Post("/", h.SetTask)          // add new task
		r.Put("/{id}", h.UpdateTask)    // update task
		r.Delete("/{id}", h.DeleteTask) // delete task
	})

	r.Post("/sendtasks", h.SendAllTasks)
}
