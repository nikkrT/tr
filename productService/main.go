package main

import (
	"context"
	"micr_course/productService/application"
	"micr_course/productService/config"
	"os/signal"
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
