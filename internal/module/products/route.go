package products

import "net/http"

// Route contains the HTTP routes for products and their variants.
type Route struct {
	handler *Handler
}

// NewRoute creates routes for products and their variants backed by handler.
func NewRoute(handler *Handler) *Route {
	return &Route{
		handler: handler,
	}
}

// Handler returns the HTTP handler for the product and variant routes.
func (r Route) Handler() http.Handler {
	mux := http.NewServeMux()

	registerProductRoutes(mux, r.handler)
	registerVariantRoutes(mux, r.handler)

	return mux
}

func registerProductRoutes(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc("GET /{$}", handler.GetMany)
	mux.HandleFunc("GET /{id}", handler.GetByID)
	mux.HandleFunc("POST /{$}", handler.Create)
	mux.HandleFunc("PUT /{id}", handler.Update)
	mux.HandleFunc("DELETE /{id}", handler.Delete)
}

func registerVariantRoutes(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc("GET /{product_id}/variants/{$}", handler.GetManyVariants)
	mux.HandleFunc("GET /{product_id}/variants/{variant_id}", handler.GetVariantByID)
	mux.HandleFunc("POST /{product_id}/variants/{$}", handler.CreateVariant)
	mux.HandleFunc("PUT /{product_id}/variants/{variant_id}", handler.UpdateVariant)
	mux.HandleFunc("DELETE /{product_id}/variants/{variant_id}", handler.DeleteVariant)
}
