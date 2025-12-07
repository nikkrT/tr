package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type Application struct {
	router http.Handler
	redis  *redis.Client
}

func New() *Application {
	app := &Application{
		redis: redis.NewClient(&redis.Options{}),
	}
	app.loadRoutes()
	return app
}

func (app *Application) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":8080",
		Handler: app.router,
	}

	err := app.redis.Ping(ctx).Err()

	if err != nil {
		return fmt.Errorf("failed to start redis: %w", err)
	}

	defer func() {
		if app.redis.Close() != nil {
			fmt.Println("redis close err:", err)
		}
	}()

	ch := make(chan error, 1)

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("application start error: %w", err)
		}
		close(ch)
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer timeoutCancel()
		return server.Shutdown(timeoutCtx)
	}
}
