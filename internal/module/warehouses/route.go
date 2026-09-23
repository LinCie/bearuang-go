package warehouses

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Route contains the HTTP routes for warehouses.
type Route struct {
	handler *Handler
}

// NewRoute creates routes for warehouses backed by handler.
func NewRoute(handler *Handler) *Route {
	return &Route{
		handler: handler,
	}
}

// Handler returns the HTTP handler for warehouse routes.
func (r Route) Handler() http.Handler {
	router := chi.NewRouter()

	router.Get("/", r.handler.GetMany)
	router.Get("/{id}", r.handler.GetByID)
	router.Post("/", r.handler.Create)
	router.Put("/{id}", r.handler.Update)
	router.Delete("/{id}", r.handler.Delete)

	return router
}
