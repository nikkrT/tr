package main

import (
	"context"
	"micr_course/orderService/applicaition"
	"micr_course/orderService/config"
	"os/signal"
	"syscall"
)

func main() {
	app := applicaition.NewApp(config.LoadConfig())

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	err := app.Start(ctx)
	if err != nil {
		panic(err)
	}
}
