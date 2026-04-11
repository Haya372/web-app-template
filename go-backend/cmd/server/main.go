package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/di"
	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/telemetry"
)

func main() {
	ctx := context.Background()

	servers, err := di.InitializeServers(ctx)
	if err \!= nil {
		panic(err)
	}

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	shutdown, err := telemetry.SetupOTelSDK(ctx)
	if err \!= nil {
		panic(err)
	}

	defer func() {
		if err := shutdown(ctx); err \!= nil {
			slog.Error("failed to shutdown telemetry", "error", err)
		}
	}()

	// Start both servers concurrently. If one fails, the shared context is
	// cancelled, which triggers a graceful shutdown of the other.
	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return servers.REST.Start(gCtx)
	})

	g.Go(func() error {
		return servers.GRPC.Start(gCtx)
	})

	if err := g.Wait(); err \!= nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}
