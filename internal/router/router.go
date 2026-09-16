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
)

type Router struct {
	router chi.Router
}

func NewRouter(db *sqlx.DB, cfg config.Config) Router {
	router := chi.NewRouter()
	router.Use(middleware.Logger, middleware.Recovery)

	authModule := auth.NewModule(db, cfg.JWTSecret)
	productCategoryModule := productcategories.NewModule(db)
	productModule := products.NewModule(db)

	router.Mount("/auth", authModule.Route.Handler())

	protected := router.With(middleware.Auth(cfg.JWTSecret))

	protected.Mount("/product-categories", productCategoryModule.Route.Handler())
	protected.Mount("/products", productModule.Route.Handler())

	return Router{
		router: router,
	}
}

func (r Router) Handler() http.Handler {
	return r.router
}
