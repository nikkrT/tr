package application

import (
	"context"
	"fmt"
	"micr_course/db"
	"micr_course/db/repo"
	"micr_course/handlers"
	"micr_course/routes"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	router http.Handler
	config Config
	db     *pgxpool.Pool
}

func NewApplication(cfg Config) *Application {
	app := &Application{
		router: nil,
		config: cfg,
		db:     nil,
	}
	return app
}

func (app *Application) Start(ctx context.Context) error {

	pool, err := db.InitDB(ctx, app.config.dbAddr)
	if err != nil {
		return fmt.Errorf("failed to init db: %v", err)
	}

	app.db = pool

	defer app.db.Close()

	productRepo := repo.NewProductRepo(app.db)
	productHandler := handlers.NewProductHandler(productRepo)
	app.router = routes.LoadRoutesProduct(productHandler)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", app.config.serverPort),
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

	select {
	case err := <-ch:
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
