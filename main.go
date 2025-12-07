package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"tr/application"
)

func main() {
	app := application.New()

	ctx, cacnel := signal.NotifyContext(context.Background(), syscall.SIGINT)

	defer cacnel()

	err := app.Start(ctx)
	if err != nil {
		fmt.Printf("application start error: %v\n", err)
	}

}
