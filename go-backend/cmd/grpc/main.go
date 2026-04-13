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
	if err := run(); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	// Initialization uses a plain context; shutdown context is signal-driven.
	shutdownCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	shutdown, err := telemetry.SetupOTelSDK(context.Background())
	if err != nil {
		return err
	}

	defer func() {
		if err := shutdown(context.Background()); err != nil {
			slog.Error("failed to shutdown telemetry", "error", err)
		}
	}()

	server, err := di.InitializeGRPCServer(context.Background())
	if err != nil {
		return err
	}

	return server.Start(shutdownCtx)
}
