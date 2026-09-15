package products

import "net/http"

// Route contains the HTTP routes for products.
type Route struct {
	handler *Handler
}

// NewRoute creates the product routes backed by handler.
func NewRoute(handler *Handler) *Route {
	return &Route{
		handler: handler,
	}
}

// Handler returns the HTTP handler for the product routes.
func (r Route) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", r.handler.GetMany)
	mux.HandleFunc("GET /{id}", r.handler.GetByID)
	mux.HandleFunc("POST /{$}", r.handler.Create)
	mux.HandleFunc("PUT /{id}", r.handler.Update)
	mux.HandleFunc("DELETE /{id}", r.handler.Delete)

	return mux
}
