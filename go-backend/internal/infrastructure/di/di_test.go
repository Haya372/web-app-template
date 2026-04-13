package di_test

// Compile-level contract tests for the di package.
//
// These tests do not execute DI initializers (which require a live DB and env
// vars). Instead they confirm that the expected symbols are exported with the
// correct signatures by assigning them to typed function variables. A missing
// or renamed symbol causes a build error, satisfying the TDD Red phase before
// the implementation exists.

import (
	"context"

	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/connectrpc"
	infrahttp "github.com/Haya372/web-app-template/go-backend/internal/infrastructure/http"

	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/di"
)

// Verify that InitializeGRPCServer is exported with the expected signature:
//
//	func(ctx context.Context) (*connectrpc.Server, error)
//
// This line will produce a compile error until the function is added to
// wire_gen.go (and wire.go) as part of the implementation task.
var _ func(context.Context) (*connectrpc.Server, error) = di.InitializeGRPCServer

// Verify that InitializeServer (REST) is still exported with its original
// signature and has not been removed or renamed.
var _ func(context.Context) (*infrahttp.Server, error) = di.InitializeServer
