package auth

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Route contains the HTTP routes for authentication.
type Route struct {
	handler *Handler
}

// NewRoute creates routes for authentication backed by handler.
func NewRoute(handler *Handler) *Route {
	return &Route{
		handler: handler,
	}
}

// Handler returns the HTTP handler for authentication routes.
func (r Route) Handler() http.Handler {
	router := chi.NewRouter()

	router.Post("/register", r.handler.Register)
	router.Post("/login", r.handler.Login)
	router.Post("/refresh", r.handler.Refresh)

	return router
}
