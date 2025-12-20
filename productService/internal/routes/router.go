package routes

import (
	"net/http"
	"productService/delivery/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func LoadRoutesProduct(productHandler *handlers.ProductHand) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
	router.Get("/bye", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bye World"))
	})

	router.Route("/products", func(r chi.Router) {
		r.Get("/{id}", productHandler.GetById)
		r.Get("/", productHandler.List)
		r.Post("/", productHandler.Create)
		r.Put("/{id}", productHandler.Update)
		r.Delete("/{id}", productHandler.Delete)
	})

	return router
}
