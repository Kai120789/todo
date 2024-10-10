package router

import (
	"net/http"
	"todo/internal/todo/middleware"

	"github.com/go-chi/chi/v5"
)

type UserRouter struct{}

type UserHandler interface {
	RegisterNewUser(w http.ResponseWriter, r *http.Request)
	AuthorizateUser(w http.ResponseWriter, r *http.Request)
	GetAuthUser(w http.ResponseWriter, r *http.Request)
	UserLogout(w http.ResponseWriter, r *http.Request)
	AddChatID(w http.ResponseWriter, r *http.Request)
}

func NewUserRouter() *UserRouter {
	return &UserRouter{}
}

func (b *UserRouter) UserRoutes(r chi.Router, h UserHandler) {
	// routes for user
	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", h.RegisterNewUser)                 // register new user
		r.Post("/login", h.AuthorizateUser)                    // login user
		r.With(middleware.JWT).Get("/", h.GetAuthUser)         // get active user, need jwt
		r.With(middleware.JWT).Delete("/logout", h.UserLogout) // logout user, need jwt
	})

	r.Post("/add-chat-id", h.AddChatID) // add chatID to table users
}
