package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Config можно вынести в отдельный пакет, как у тебя в других сервисах
type Config struct {
	Port                string
	ProductServiceURL   string
	OrderServiceURL     string
	NotificationService string
}

func main() {
	// В идеале загружать это через твой пакет config
	cfg := Config{
		Port:              "8084",
		ProductServiceURL: "http://localhost:8081", // Укажи здесь порты из docker-compose
		OrderServiceURL:   "http://localhost:8082",
	}

	router := chi.NewRouter()

	// Полезные мидлвари в стиле твоего кода
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Настраиваем прокси для сервисов
	productProxy := createReverseProxy(cfg.ProductServiceURL)
	orderProxy := createReverseProxy(cfg.OrderServiceURL)

	// Маршрутизация запросов к микросервисам
	// Все запросы, начинающиеся с /api/v1/products будут перенаправляться в ProductService
	router.Route("/api/v1/products", func(r chi.Router) {
		r.HandleFunc("/*", func(w http.ResponseWriter, req *http.Request) {
			// Убираем префикс /api/v1 (если ProductService его не ждет)
			// Если ProductService ждет пути вроде /products, то нужно использовать http.StripPrefix
			productProxy.ServeHTTP(w, req)
		})
	})

	// Все запросы, начинающиеся с /api/v1/orders будут перенаправляться в OrderService
	router.Route("/api/v1/orders", func(r chi.Router) {
		r.HandleFunc("/*", func(w http.ResponseWriter, req *http.Request) {
			orderProxy.ServeHTTP(w, req)
		})
	})

	// Сервер Gateway
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful Shutdown (по твоему стилю)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		log.Printf("Gateway is starting on port %s...", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down Gateway...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Gateway shutdown failed: %v", err)
	}

	log.Println("Gateway stopped gracefully")
}

// createReverseProxy создает httputil.ReverseProxy для целевого URL
func createReverseProxy(target string) *httputil.ReverseProxy {
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Invalid target URL %s: %v", target, err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Опционально: можно изменять запросы перед отправкой в микросервис
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		// Если нужно передать какие-то хидеры от gateway до сервисов, делаем это здесь
		req.Header.Set("X-Gateway-Proxy", "true")
	}

	return proxy
}
