package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"productService/internal/model"
	"productService/internal/service"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ProductHandler interface {
	CreateProduct(ctx context.Context, product model.Product) error
	ReadProduct(ctx context.Context, id int) (model.Product, error)
	ReadAll(ctx context.Context, filteredBy string) ([]model.Product, error)
	UpdateProduct(ctx context.Context, product model.Product) error
	DeleteProduct(ctx context.Context, id int) error
}

type ProductHand struct {
	Service ProductHandler
}

func NewProductHandler(r *service.ProductService) *ProductHand {
	return &ProductHand{
		Service: r,
	}
}

func (p *ProductHand) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Price       int    `json:"price"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, fmt.Sprintf("decode error: %v", err), http.StatusBadRequest)
		return
	}
	product := model.Product{
		Name:        body.Name,
		Description: body.Description,
		Price:       body.Price,
	}
	err := p.Service.CreateProduct(r.Context(), product)
	if err != nil {
		http.Error(w, fmt.Sprintf("insert error: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (p *ProductHand) GetById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, _ := strconv.Atoi(id)
	product, err := p.Service.ReadProduct(r.Context(), idInt)
	if err != nil {
		if err == service.ErrNotFound {
			http.Error(w, fmt.Sprintf("Product with id %v not found", idInt), http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("get product by id %v error: %v", idInt, err), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(product); err != nil {
		http.Error(w, "Failed to encode product", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (p *ProductHand) List(w http.ResponseWriter, r *http.Request) {
	filteredBy := r.URL.Query().Get("filteredBy")
	fmt.Println(filteredBy)
	products := []model.Product{}
	products, err := p.Service.ReadAll(r.Context(), filteredBy)
	if err != nil {
		http.Error(w, fmt.Sprintf("list products error: %v", err), http.StatusInternalServerError)
	}
	if err := json.NewEncoder(w).Encode(products); err != nil {
		http.Error(w, "Failed to encode products", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (p *ProductHand) Update(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Id    int `json:"id"`
		Price int `json:"price"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, fmt.Sprintf("decode error: %v", err), http.StatusBadRequest)
		return
	}
	product, err := p.Service.ReadProduct(r.Context(), body.Id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Product with id %v not found", body.Id), http.StatusNotFound)
		return
	}
	product.Price = body.Price
	err = p.Service.UpdateProduct(r.Context(), product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}
	if err := json.NewEncoder(w).Encode(product); err != nil {
		http.Error(w, "Failed to encode product", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *ProductHand) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, _ := strconv.Atoi(id)
	err := h.Service.DeleteProduct(r.Context(), idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
