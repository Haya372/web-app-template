package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/di"
	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/telemetry"
)

func main() {
	ctx := context.Background()

	server, err := di.InitializeServer(ctx)
	if err != nil {
		panic(err)
	}

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	shutdown, err := telemetry.SetupOTelSDK(ctx)
	if err != nil {
		panic(err)
	}

	defer func() {
		err := shutdown(ctx)
		if err != nil {
			slog.Error("failed to shutdown telemetry", "error", err)
		}
	}()

	if err := server.Start(ctx); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}
