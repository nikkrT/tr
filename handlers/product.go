package handlers

import (
	"encoding/json"
	"fmt"
	"micr_course/db/repo"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Product struct {
	Repo *repo.ProductRepo
}

func NewProductHandler(r *repo.ProductRepo) *Product {
	return &Product{
		Repo: r,
	}
}

func (p *Product) GetById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	fmt.Println(id + "ewfewfewfewf")
	idInt, _ := strconv.Atoi(id)
	product, err := p.Repo.FindById(r.Context(), idInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	if err := json.NewEncoder(w).Encode(product); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
