package application

import (
	"context"
	"fmt"
	"productService/config"
	"productService/delivery/events"

	"productService/delivery/grpc"
	"productService/delivery/handlers"

	"productService/internal/db"
	"productService/internal/db/repo"
	"productService/internal/routes"
	"productService/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
	g "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"net"
	"net/http"
	pb "productService/pkg/proto"
	"time"
)

type Application struct {
	router   http.Handler
	config   *config.Config
	db       *pgxpool.Pool
	rabbitmq *events.RabbitMQPublisher
}

func NewApplication(cfg *config.Config) *Application {
	app := &Application{
		router: nil,
		config: cfg,
		db:     nil,
	}
	return app
}

func (app *Application) Start(ctx context.Context) error {

	pool, err := db.InitDB(ctx, app.config.Postgres)
	if err != nil {
		return fmt.Errorf("failed to init db: %v", err)
	}

	app.db = pool

	defer app.db.Close()

	productRepo := repo.NewProductRepo(app.db)

	app.rabbitmq, err = events.NewProducerSetup(app.config)

	if err != nil {
		return fmt.Errorf("failed to init rabbitmq: %v", err)
	}

	productService := service.NewProductService(productRepo, app.rabbitmq)

	productHandler := handlers.NewProductHandler(productService)
	grpcServer := grpc.NewGRPCServer(productService)

	app.router = routes.LoadRoutesProduct(productHandler)

	defer app.rabbitmq.Close()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", app.config.Server.ServerPort),
		Handler:      app.router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}

	ch := make(chan error)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch)
	}()
	go func() {
		lis, err := net.Listen("tcp", ":8081")
		if err != nil {
			ch <- fmt.Errorf("failed to listen: %v", err)
		}
		s := g.NewServer()
		reflection.Register(s)
		pb.RegisterProductServiceServer(s, grpcServer) // Регистрируем нашу реализацию
		ch <- s.Serve(lis)
	}()

	select {
	case err := <-ch:
		fmt.Println(err)
		return err
	case <-ctx.Done():
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		fmt.Println("\napplication shutdown started")
		defer cancel()
		if err := server.Shutdown(timeoutCtx); err != nil {
			return err
		}
	}
	return nil
}
