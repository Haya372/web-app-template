package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Haya372/web-app-template/backend/internal/infrastructure/di"
)

func main() {
	ctx := context.Background()

	server, err := di.InitializeServer(ctx)
	if err != nil {
		panic(err)
	}

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := server.Start(ctx); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}
