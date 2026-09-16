package productcategories

import "net/http"

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
	mux := http.NewServeMux()

	registerProductCategoryRoutes(mux, r.handler)

	return mux
}

func registerProductCategoryRoutes(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc("GET /{$}", handler.GetMany)
	mux.HandleFunc("GET /{id}", handler.GetByID)
	mux.HandleFunc("POST /{$}", handler.Create)
	mux.HandleFunc("PUT /{id}", handler.Update)
	mux.HandleFunc("DELETE /{id}", handler.Delete)
}
