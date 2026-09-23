package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"bearuang-go/internal/config"
	"bearuang-go/internal/middleware"
	"bearuang-go/internal/module/auth"
	"bearuang-go/internal/module/productcategories"
	"bearuang-go/internal/module/products"
	"bearuang-go/internal/module/warehouses"
)

type Router struct {
	router chi.Router
}

func NewRouter(db *sqlx.DB, cfg config.Config) Router {
	router := chi.NewRouter()
	router.Use(middleware.Logger, middleware.Recovery)

	authModule := auth.NewModule(db, cfg)
	productCategoryModule := productcategories.NewModule(db)
	productModule := products.NewModule(db)
	warehouseModule := warehouses.NewModule(db)

	router.Mount("/auth", authModule.Route.Handler())

	protected := router.With(middleware.Auth(cfg))

	protected.Mount("/product-categories", productCategoryModule.Route.Handler())
	protected.Mount("/products", productModule.Route.Handler())
	protected.Mount("/warehouses", warehouseModule.Route.Handler())

	return Router{
		router: router,
	}
}

func (r Router) Handler() http.Handler {
	return r.router
}
