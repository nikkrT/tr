package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	pbProduct "micr_course/pkg/proto/productService"
)

// Request структура для JSON (то, что приходит от фронтенда)
type CreateProductHTTPRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       uint32 `json:"price"`
}

type ProductHandler struct {
	client pbProduct.ProductServiceClient
}

// Конструктор
func NewProductHandler(client pbProduct.ProductServiceClient) *ProductHandler {
	return &ProductHandler{
		client: client,
	}
}

// HTTP хендлер для создания товара
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	// 1. Декодируем JSON из HTTP запроса
	var reqBody CreateProductHTTPRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 2. Формируем gRPC запрос
	grpcReq := &pbProduct.CreateProductRequest{
		Name:        reqBody.Name,
		Description: reqBody.Description,
		Price:       reqBody.Price,
	}

	// Используем контекст с таймаутом, чтобы gateway не зависал, если productService недоступен
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// 3. Делаем вызов по gRPC к productService
	grpcResp, err := h.client.CreateProduct(ctx, grpcReq)
	if err != nil {
		http.Error(w, "Failed to create product via gRPC: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Отправляем успешный HTTP ответ (перекладываем данные из gRPC ответа в JSON)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": grpcResp.Id,
	})
}
