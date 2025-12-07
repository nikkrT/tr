package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"tr/model"
	"tr/repository/order"

	"github.com/google/uuid"
)

type Order struct {
	Repo *order.RedisRepository
}

func (h *Order) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CustomerId uuid.UUID        `json:"customer_id"`
		LineItems  []model.LineItem `json:"line_items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	now := time.Now()
	order := model.Order{
		OrderId:    rand.Uint64(),
		CustomerId: body.CustomerId,
		LineItems:  body.LineItems,
		CreatedAt:  &now,
	}

	err := h.Repo.Insert(r.Context(), order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	res, err := json.Marshal(order)
	if err != nil {
		fmt.Println("failed to marshal order")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(res)
	w.WriteHeader(http.StatusCreated)
}

func (*Order) List(w http.ResponseWriter, r *http.Request) {
	fmt.Println("List")
}

func (*Order) GetById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetById")
}
func (*Order) UpdateById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UpdateById")
}

func (*Order) DeleteById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UpdateById")
}
