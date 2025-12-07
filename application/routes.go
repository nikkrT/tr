package application

import (
	"net/http"
	"tr/handlers"
	"tr/repository/order"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *Application) loadRoutes() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/order", a.loadOrderRoutes)
	a.router = router
}

func (a *Application) loadOrderRoutes(router chi.Router) {
	orderHandler := &handlers.Order{
		Repo: &order.RedisRepository{
			Client: a.redis,
		},
	}

	router.Post("/", orderHandler.Create)
	router.Get("/", orderHandler.List)
	router.Get("/{id}", orderHandler.GetById)
	router.Delete("/{id}", orderHandler.DeleteById)
	router.Put("/{id}", orderHandler.UpdateById)
}
