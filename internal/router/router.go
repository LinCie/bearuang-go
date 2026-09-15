package router

import (
	"net/http"

	"github.com/jmoiron/sqlx"

	"bearuang-go/internal/httpx"
	"bearuang-go/internal/middleware"
)

type Router struct {
	mux *http.ServeMux
}

func NewRouter(db *sqlx.DB) Router {
	mux := http.NewServeMux()

	return Router{
		mux: mux,
	}
}

func (r Router) Handler() http.Handler {
	return httpx.NewChain(
		middleware.Logger,
		middleware.Recovery,
	).Then(r.mux)
}
