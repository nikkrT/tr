package application

import (
	""
	"context"
	"fmt"
	"net/http"
)

type Application struct {
	router http.Handler
}

func New() *Application {
	app := &Application{
		router: loadRouters(),
	}
	return app
}

func (app *Application) Start(ctx context.Context) error {
	server := http.Server{
		Addr:    ":8080",
		Handler: app.router,
	}

	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("application start error: %w", err)
	}

	return nil
}
