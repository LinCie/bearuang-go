package inventory

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Route contains the HTTP routes for warehouses and inventory.
type Route struct {
	handler *Handler
}

// NewRoute creates routes for warehouses and inventory backed by handler.
func NewRoute(handler *Handler) *Route {
	return &Route{
		handler: handler,
	}
}

// Handler returns the HTTP handler for warehouse and inventory routes.
func (r Route) Handler() http.Handler {
	router := chi.NewRouter()

	registerWarehouseRoutes(router, r.handler)
	registerStockRoutes(router, r.handler)

	return router
}

func registerWarehouseRoutes(router chi.Router, handler *Handler) {
	router.Get("/", handler.GetManyWarehouses)
	router.Post("/", handler.CreateWarehouse)
	router.Get("/{id}", handler.GetWarehouseByID)
	router.Put("/{id}", handler.UpdateWarehouse)
	router.Delete("/{id}", handler.DeleteWarehouse)
}

func registerStockRoutes(router chi.Router, handler *Handler) {
	router.Get("/{warehouse_id}/stocks/", handler.GetWarehouseStocks)
	router.Get("/{warehouse_id}/stocks/{variant_id}", handler.GetWarehouseStock)
	router.Get("/{warehouse_id}/adjustments/", handler.GetStockAdjustments)
	router.Post("/{warehouse_id}/adjustments/", handler.CreateAdjustment)
}
