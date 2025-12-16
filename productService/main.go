package main

import (
	"context"
	"os/signal"
	"productService/application"
	"syscall"
)

func main() {
	app := application.NewApplication(application.LoadConfig())

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	err := app.Start(ctx)
	if err != nil {
		panic(err)
	}
}
