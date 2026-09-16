package auth

import "github.com/jmoiron/sqlx"

// Module wires the authentication repository, service, handler, and routes.
type Module struct {
	Route *Route
}

// NewModule creates the authentication module backed by db.
func NewModule(db *sqlx.DB, jwtSecret string) Module {
	repository := NewPostgresRepository(db)
	service := NewService(repository, jwtSecret)
	handler := NewHandler(service)

	return Module{
		Route: NewRoute(handler),
	}
}
