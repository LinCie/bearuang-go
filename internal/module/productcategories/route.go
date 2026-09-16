package productcategories

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Route contains the HTTP routes for product categories.
type Route struct {
	handler *Handler
}

// NewRoute creates routes for product categories backed by handler.
func NewRoute(handler *Handler) *Route {
	return &Route{
		handler: handler,
	}
}

// Handler returns the HTTP handler for the product category routes.
func (r Route) Handler() http.Handler {
	router := chi.NewRouter()

	registerProductCategoryRoutes(router, r.handler)

	return router
}

func registerProductCategoryRoutes(router chi.Router, handler *Handler) {
	router.Get("/", handler.GetMany)
	router.Get("/{id}", handler.GetByID)
	router.Post("/", handler.Create)
	router.Put("/{id}", handler.Update)
	router.Delete("/{id}", handler.Delete)
}
