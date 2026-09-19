package auth

import (
	"github.com/jmoiron/sqlx"

	"bearuang-go/internal/config"
)

// Module wires the authentication repository, service, handler, and routes.
type Module struct {
	Route *Route
}

// NewModule creates the authentication module backed by db.
func NewModule(db *sqlx.DB, cfg config.Config) Module {
	repository := NewPostgresRepository(db)
	service := NewService(repository, cfg.JWTSecret, NewArgon2Password())
	handler := NewHandler(service)

	return Module{
		Route: NewRoute(handler),
	}
}
