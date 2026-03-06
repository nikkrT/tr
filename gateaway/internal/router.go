package internal

import (
	"micr_course/productService/product/interfaces/net_requests"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func LoadRoutesProduct(productHandler *net_requests.ProductHand) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/hello", func(w net_requests.ResponseWriter, r *net_requests.Request) {
		w.Write([]byte("Hello World"))
	})
	router.Get("/bye", func(w net_requests.ResponseWriter, r *net_requests.Request) {
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
