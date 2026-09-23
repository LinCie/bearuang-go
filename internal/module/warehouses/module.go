package warehouses

import "github.com/jmoiron/sqlx"

// Module wires the warehouses repository, service, handler, and routes.
type Module struct {
	Route *Route
}

// NewModule creates the warehouses module backed by db.
func NewModule(db *sqlx.DB) Module {
	repository := NewPostgresRepository(db)
	service := NewService(repository)
	handler := NewHandler(service)

	return Module{
		Route: NewRoute(handler),
	}
}
