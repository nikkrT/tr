package main

import (
	"context"
	"fmt"
	"notificationService/application"
	"notificationService/config"
	"os/signal"
	"syscall"
)

func main() {
	app := application.NewApplication(config.LoadConfig())
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	err := app.Start(ctx)
	if err != nil {
		fmt.Printf("Failed to start application: %s", err)
	}
}
