package router

import (
	"net/http"
	"todo/internal/todo/middleware"

	"github.com/go-chi/chi/v5"
)

type BoardsRouter struct{}

type BoardsHandler interface {
	SetBoard(w http.ResponseWriter, r *http.Request)
	GetAllBoards(w http.ResponseWriter, r *http.Request)
	GetBoard(w http.ResponseWriter, r *http.Request)
	UpdateBoard(w http.ResponseWriter, r *http.Request)
	DeleteBoard(w http.ResponseWriter, r *http.Request)
	User2Board(w http.ResponseWriter, r *http.Request)
}

func NewBoardsRouter() *BoardsRouter {
	return &BoardsRouter{}
}

func (b *BoardsRouter) BoardsRoutes(r chi.Router, h BoardsHandler) {
	// Routes for boards
	r.Route("/api/boards", func(r chi.Router) {
		r.Use(middleware.JWT)            // need jwt for all methods
		r.Get("/", h.GetAllBoards)       // get all boards
		r.Get("/{id}", h.GetBoard)       // get board with id
		r.Post("/", h.SetBoard)          // add new board
		r.Put("/{id}", h.UpdateBoard)    // update board
		r.Delete("/{id}", h.DeleteBoard) // delete board
		r.Post("/{id}", h.User2Board)    // add user to board
	})
}
