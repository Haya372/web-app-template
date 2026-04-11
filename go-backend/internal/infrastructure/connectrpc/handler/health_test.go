package handler_test

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/connectrpc/handler"
	v1 "github.com/Haya372/web-app-template/go-backend/internal/infrastructure/grpc/gen/health/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthHandler_Check_ReturnsServing(t *testing.T) {
	tests := []struct {
		name string
		ctx  func() context.Context
	}{
		{
			name: "returns serving status with no error",
			ctx:  context.Background,
		},
		{
			name: "returns serving status even with cancelled context",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()

				return ctx
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler.NewHealthHandler()
			req := connect.NewRequest(&v1.CheckRequest{})

			resp, err := h.Check(tt.ctx(), req)

			require.NoError(t, err)
			assert.Equal(t, v1.ServingStatus_SERVING_STATUS_SERVING, resp.Msg.GetStatus())
		})
	}
}
