package router

import (
	"net/http"
	"todo/internal/todo/middleware"

	"github.com/go-chi/chi/v5"
)

type StatusesRouter struct{}

type StatusesHandler interface {
	SetStatus(w http.ResponseWriter, r *http.Request)
	DeleteStatus(w http.ResponseWriter, r *http.Request)
	HelloWorld(w http.ResponseWriter, r *http.Request)
}

func NewStatusesRouter() *StatusesRouter {
	return &StatusesRouter{}
}

func (b *StatusesRouter) StatusesRoutes(r chi.Router, h StatusesHandler) {
	// Routes for statuses
	r.Route("/api/status", func(r chi.Router) {
		r.Use(middleware.JWT)         // need jwt for all methods
		r.Post("/", h.SetStatus)      // add new status
		r.Delete("/", h.DeleteStatus) // delete status
	})

	r.Get("/", h.HelloWorld)
}
