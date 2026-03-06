package application

import (
	"context"
	"fmt"
	pb "micr_course/pkg/proto/productService"
	"micr_course/productService/config"
	"micr_course/productService/product/infrastructure/messaging"
	"micr_course/productService/product/infrastructure/postgres"
	"micr_course/productService/product/infrastructure/postgres/repo"
	myGrpc "micr_course/productService/product/interfaces/grpc"
	"micr_course/productService/product/interfaces/net_requests"
	"micr_course/productService/product/interfaces/net_requests/routes"
	"micr_course/productService/product/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"net"
	"net/http"
	"time"
)

type Application struct {
	router   http.Handler
	config   *config.Config
	db       *pgxpool.Pool
	rabbitmq *messaging.RabbitMQPublisher
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

	pool, err := postgres.InitDB(ctx, app.config.Postgres)
	if err != nil {
		return fmt.Errorf("failed to init postgres: %v", err)
	}

	app.db = pool

	defer app.db.Close()

	productRepo := repo.NewProductRepo(app.db)

	app.rabbitmq, err = messaging.NewProducerSetup(&app.config.RabbitMQ)

	if err != nil {
		return fmt.Errorf("failed to init rabbitmq: %v", err)
	}

	productService := service.NewProductService(productRepo, app.rabbitmq)

	productHandler := net_requests.NewProductHandler(productService)
	grpcServer := myGrpc.NewGRPCServer(productService)

	app.router = routes.LoadRoutesProduct(productHandler)
	defer app.rabbitmq.Close()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", app.config.Server.ServerPort),
		Handler:      app.router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}

	ch := make(chan error, 2)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch)
	}()
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", app.config.GRPC.Port))
		if err != nil {
			ch <- fmt.Errorf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		reflection.Register(s)
		pb.RegisterProductServiceServer(s, grpcServer)
		ch <- s.Serve(lis)
		close(ch)
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
