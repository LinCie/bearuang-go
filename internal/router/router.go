package router

import (
	"net/http"

	"github.com/jmoiron/sqlx"

	"bearuang-go/internal/httpx"
	"bearuang-go/internal/middleware"
	"bearuang-go/internal/module/products"
)

type Router struct {
	mux *http.ServeMux
}

func NewRouter(db *sqlx.DB) Router {
	mux := http.NewServeMux()
	productModule := products.NewModule(db)

	mux.Handle(
		"/products/",
		http.StripPrefix("/products", productModule.Route.Handler()),
	)

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
