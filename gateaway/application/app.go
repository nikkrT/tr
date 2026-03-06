package application

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"micr_course/gateaway/config"
)

type Application struct {
	router *http.Handler
	cfg *config.Config
}

func NewApplication(cfg config.Config) *Application {
	app:= &Application{
		cfg: &cfg,
	}
	app.router := LoadRoutes()
	return app
}
