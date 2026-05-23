// @title Assistant API
// @version 1.0
// @description Example Swagger setup for Golang Gin.
// @schemes http
// @host 127.0.0.1:9876
// @BasePath /
package main

import (
	"assistant/internal/bootstrap"
	"context"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := bootstrap.Run(ctx); err != nil {
		log.Fatalf("server exited with error: %v", err)
	}
}
