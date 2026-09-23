package inventory

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Route contains the HTTP routes for inventory.
type Route struct {
	handler *Handler
}

// NewRoute creates inventory routes backed by handler.
func NewRoute(handler *Handler) *Route {
	return &Route{handler: handler}
}

// Handler returns the HTTP handler for inventory routes.
func (r Route) Handler() http.Handler {
	router := chi.NewRouter()

	router.Get("/balances/", r.handler.GetBalances)
	router.Get("/balances/{warehouse_id}/{variant_id}", r.handler.GetBalance)
	router.Post("/adjustments/", r.handler.Adjust)
	router.Get("/movements/", r.handler.GetMovements)
	router.Get("/movements/{id}", r.handler.GetMovementByID)

	return router
}
