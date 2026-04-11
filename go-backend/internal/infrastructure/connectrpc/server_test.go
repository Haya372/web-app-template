package connectrpc_test

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/connectrpc"
	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/connectrpc/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewServer_UsesEnvPort verifies that NewServer respects APP_GRPC_PORT.
func TestNewServer_UsesEnvPort(t *testing.T) {
	t.Setenv("APP_GRPC_PORT", "19090")

	s := connectrpc.NewServer(handler.NewHealthHandler())

	assert.Equal(t, ":19090", s.Addr())
}

// TestNewServer_DefaultPort verifies that NewServer falls back to port 8081
// when APP_GRPC_PORT is not set.
func TestNewServer_DefaultPort(t *testing.T) {
	t.Setenv("APP_GRPC_PORT", "")

	s := connectrpc.NewServer(handler.NewHealthHandler())

	assert.Equal(t, ":8081", s.Addr())
}

// TestServer_Start_ShutdownOnContextCancel verifies that Start returns without
// error when the context is cancelled, confirming graceful shutdown behaviour.
func TestServer_Start_ShutdownOnContextCancel(t *testing.T) {
	t.Setenv("APP_GRPC_PORT", "18081")

	h := handler.NewHealthHandler()
	s := connectrpc.NewServer(h)

	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error, 1)

	go func() {
		errCh <- s.Start(ctx)
	}()

	// Give the server a moment to bind before cancelling.
	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case err := <-errCh:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("Start did not return within timeout after context cancel")
	}
}

// TestServer_Start_AlreadyCancelledContext verifies that Start returns immediately
// when the context is already cancelled before Start is called.
func TestServer_Start_AlreadyCancelledContext(t *testing.T) {
	t.Setenv("APP_GRPC_PORT", "18082")

	h := handler.NewHealthHandler()
	s := connectrpc.NewServer(h)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	errCh := make(chan error, 1)

	go func() {
		errCh <- s.Start(ctx)
	}()

	select {
	case err := <-errCh:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("Start did not return within timeout for pre-cancelled context")
	}
}

// TestServer_Start_PortAlreadyBound verifies that Start returns a non-nil error
// when the configured port is already in use.
// The listener must bind on all interfaces (":0") to conflict with the server's
// all-interface bind (":PORT"). A loopback-only bind (127.0.0.1:PORT) does not
// conflict on macOS.
func TestServer_Start_PortAlreadyBound(t *testing.T) {
	// Bind on all interfaces with an OS-assigned random port.
	ln, err := net.Listen("tcp", ":0") //nolint:gosec // test-only ephemeral listener
	require.NoError(t, err)

	defer ln.Close()

	tcpAddr, ok := ln.Addr().(*net.TCPAddr)
	require.True(t, ok, "expected *net.TCPAddr")

	port := strconv.Itoa(tcpAddr.Port)
	t.Setenv("APP_GRPC_PORT", port)

	h := handler.NewHealthHandler()
	s := connectrpc.NewServer(h)

	err = s.Start(context.Background())
	assert.Error(t, err)
}
