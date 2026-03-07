package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	pbOrder "micr_course/pkg/proto/orderService"
)

// Request структуры для JSON (то, что приходит от фронтенда)
type CreateOrderHTTPRequest struct {
	ProductID uint32 `json:"product_id"`
	Quantity  uint32 `json:"quantity"`
}

type OrderHandler struct {
	client pbOrder.OrderServiceClient
}

// Конструктор
func NewOrderHandler(client pbOrder.OrderServiceClient) *OrderHandler {
	return &OrderHandler{
		client: client,
	}
}

// HTTP хендлер для создания заказа
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	// 1. Декодируем JSON из HTTP запроса
	var reqBody CreateOrderHTTPRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 2. Формируем gRPC запрос
	// (Ожидается, что ты добавишь CreateOrderRequest в свой proto файл)
	grpcReq := &pbOrder.OrderRequest{
		Id:       reqBody.ProductID,
		Quantity: reqBody.Quantity,
	}

	// Используем контекст с таймаутом, чтобы gateway не зависал, если orderService недоступен
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// 3. Делаем вызов по gRPC к orderService!
	grpcResp, err := h.client.CreateOrder(ctx, grpcReq)
	if err != nil {
		// В реальном проекте тут лучше разбирать коды ошибок gRPC (codes.NotFound, codes.Internal и тд)
		http.Error(w, "Failed to create order via gRPC: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Отправляем успешный HTTP ответ (перекладываем данные из gRPC ответа в JSON)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	var res = "Order created successfully"
	if grpcResp.OrderId == 0 {
		res = "Order created unsuccessfully"
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  res,
		"order_id": grpcResp.OrderId,
		"success":  grpcResp.Success,
	})
}
