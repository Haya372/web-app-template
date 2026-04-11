package connectrpc

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/connectrpc/handler"
	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/grpc/gen/health/v1/healthv1connect"
)

const (
	gracefulTimeout   = 5 * time.Second
	readHeaderTimeout = 10 * time.Second
)

// Server wraps a net/http.Server that serves Connect-RPC endpoints.
type Server struct {
	httpServer *http.Server
}

// NewServer constructs a Connect-RPC Server listening on APP_GRPC_PORT (default :8081).
func NewServer(h *handler.HealthHandler) *Server {
	port := os.Getenv("APP_GRPC_PORT")
	if port == "" {
		port = "8081"
	}

	mux := http.NewServeMux()

	path, connectHandler := healthv1connect.NewHealthServiceHandler(h)
	mux.Handle(path, connectHandler)

	return &Server{
		httpServer: &http.Server{
			Addr:              ":" + port,
			Handler:           otelhttp.NewHandler(mux, "connectrpc"),
			ReadHeaderTimeout: readHeaderTimeout,
		},
	}
}

// Addr returns the configured listen address (e.g. ":8081").
func (s *Server) Addr() string {
	return s.httpServer.Addr
}

// Start begins serving and blocks until ctx is cancelled, then performs a
// graceful shutdown. A ListenAndServe failure is returned directly.
func (s *Server) Start(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}

		close(errCh)
	}()

	select {
	case <-ctx.Done():
		// Use WithoutCancel so the shutdown context is not immediately cancelled
		// by the already-done parent ctx.
		shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), gracefulTimeout)
		defer cancel()

		shutdownErr := s.httpServer.Shutdown(shutdownCtx)

		if listenErr := <-errCh; listenErr != nil {
			return listenErr
		}

		return shutdownErr
	case err := <-errCh:
		return err
	}
}
