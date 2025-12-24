package main

import (
	"context"
	"os/signal"
	"productService/application"
	"productService/config"
	"syscall"
)

func main() {
	app := application.NewApplication(config.LoadConfig())

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	err := app.Start(ctx)
	if err != nil {
		panic(err)
	}
}
