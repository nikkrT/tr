package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"micr_course/handlers"
	"net/http"
)

func LoadRoutesProduct(productHandler *handlers.Product) *chi.Mux {
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
	})

	return router
}
