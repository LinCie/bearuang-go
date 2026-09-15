package products

import "github.com/jmoiron/sqlx"

// Module wires the products repository, service, handler, and routes.
type Module struct {
	Route *Route
}

// NewModule creates the products module backed by db.
func NewModule(db *sqlx.DB) Module {
	repository := NewPostgresDatabase(db)
	service := NewService(repository)
	handler := NewHandler(service)

	return Module{
		Route: NewRoute(handler),
	}
}
