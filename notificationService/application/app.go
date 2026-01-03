package application

import (
	"context"
	"fmt"
	"net/http"
	"notificationService/config"
	"notificationService/delivery"
	"notificationService/internal/infrastucture"
	"notificationService/internal/usecase"

	"time"
)

type Application struct {
	config *config.Config
}

func NewApplication(config *config.Config) *Application {
	app := &Application{
		config: config,
	}
	return app
}

func (app *Application) Start(ctx context.Context) error {

	infra := infrastucture.NewInfrastucture()

	service := usecase.NewService(infra)

	consumer, err := delivery.NewRabbitMQConsumer(service, app.config.Amqp)
	if err != nil {
		return fmt.Errorf("delivery.NewRabbitMQConsumer: %w", err)
	}

	server := &http.Server{
		Addr: fmt.Sprintf(":%s", app.config.Server.Port),
	}

	ch := make(chan error)

	go func() {
		err = consumer.StartConsumers(app.config.Amqp.WorkerPoolSize, app.config.Amqp.ConsumerTag)
		if err != nil {
			ch <- fmt.Errorf("consumer.StartConsumers: %w", err)
		}
		close(ch)
	}()

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ch <- fmt.Errorf("Failed to listen on port %d: %s", app.config.Server.Port, err)
		}
		close(ch)
	}()

	select {
	case err := <-ch:
		return fmt.Errorf("Received error: %s", err)
	case <-ctx.Done():
		timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer timeoutCancel()
		fmt.Println("\nReceived shutdown signal, shutting down gracefully")
		if consumer != nil {
			consumer.Close()
		}
		if err := server.Shutdown(timeoutCtx); err != nil {
			return fmt.Errorf("Failed to shutdown: %s", err)
		}
	}
	return nil
}
