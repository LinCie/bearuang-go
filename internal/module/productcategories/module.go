package productcategories

import "github.com/jmoiron/sqlx"

// Module wires the product categories repository, service, handler, and routes.
type Module struct {
	Route *Route
}

// NewModule creates the product categories module backed by db.
func NewModule(db *sqlx.DB) Module {
	repository := NewPostgresRepository(db)
	service := NewService(repository)
	handler := NewHandler(service)

	return Module{
		Route: NewRoute(handler),
	}
}
