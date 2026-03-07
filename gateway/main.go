package main

import (
	"context"
	"fmt"
	"log"
	"micr_course/gateway/config"
	"micr_course/gateway/handlers"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// Импортируем твои сгенерированные proto-файлы
	pbOrder "micr_course/pkg/proto/orderService"
	pbProduct "micr_course/pkg/proto/productService"
)

func main() {
	// 1. Устанавливаем gRPC соединение с orderService
	// Не забудь использовать insecure.NewCredentials(), как мы обсуждали ранее!
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	orderGrpcAddr := cfg.OrderServiceURL

	conn, err := grpc.NewClient(
		orderGrpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to order gRPC service: %v", err)
	}
	defer conn.Close()

	// Создаем gRPC клиента для orderService
	orderClient := pbOrder.NewOrderServiceClient(conn)

	productGrpcAddr := cfg.ProductServiceURL
	productConn, err := grpc.NewClient(
		productGrpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to product gRPC service: %v", err)
	}
	defer productConn.Close()

	// Создаем gRPC клиента для productService
	productClient := pbProduct.NewProductServiceClient(productConn)

	// 2. Инициализируем наши HTTP обработчики (передаем им gRPC клиента)
	orderHandler := handlers.NewOrderHandler(orderClient)
	productHandler := handlers.NewProductHandler(productClient)

	// 3. Настраиваем роутер chi
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Регистрируем эндпоинт для создания заказа
	router.Route("/api/v1/orders", func(r chi.Router) {
		r.Post("/", orderHandler.CreateOrder) // POST /api/v1/orders
	})

	// Регистрируем эндпоинт для создания товара
	router.Route("/api/v1/products", func(r chi.Router) {
		r.Post("/", productHandler.CreateProduct) // POST /api/v1/products
	})

	// 4. Запускаем HTTP сервер Gateway
	port := cfg.Port
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	// Graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		log.Printf("Gateway HTTP server started on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down Gateway...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Gateway shutdown failed: %v", err)
	}
}
