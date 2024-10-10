package router

import (
	"net/http"
	"todo/internal/todo/transport/http/handler"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	Boards   BoardsRouter
	Statuses StatusesRouter
	Tasks    TasksRouter
	User     UserRouter
}

func New(h *handler.TodoHandler) http.Handler {
	r := chi.NewRouter()

	router := &Router{
		Boards:   *NewBoardsRouter(),
		Statuses: *NewStatusesRouter(),
		Tasks:    *NewTasksRouter(),
		User:     *NewUserRouter(),
	}

	router.Boards.BoardsRoutes(r, &h.BoardsHandler)
	router.Statuses.StatusesRoutes(r, &h.StatusesHandler)
	router.Tasks.TasksRoutes(r, &h.TasksHandler)
	router.User.UserRoutes(r, &h.UserHandler)

	return r
}
