package router

import (
	"net/http"

	"github.com/jmoiron/sqlx"

	"bearuang-go/internal/httpx"
	"bearuang-go/internal/middleware"
	"bearuang-go/internal/module/productcategories"
	"bearuang-go/internal/module/products"
)

type Router struct {
	mux *http.ServeMux
}

func NewRouter(db *sqlx.DB, jwtSecret string) Router {
	mux := http.NewServeMux()
	productCategoryModule := productcategories.NewModule(db)
	productModule := products.NewModule(db)
	protected := httpx.NewChain(middleware.Auth(jwtSecret))

	mux.Handle(
		"/product-categories/",
		protected.Then(
			http.StripPrefix("/product-categories", productCategoryModule.Route.Handler()),
		),
	)

	mux.Handle(
		"/products/",
		protected.Then(
			http.StripPrefix("/products", productModule.Route.Handler()),
		),
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
