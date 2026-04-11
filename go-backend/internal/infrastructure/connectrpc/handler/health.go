package handler

import (
	"context"

	"connectrpc.com/connect"
	v1 "github.com/Haya372/web-app-template/go-backend/internal/infrastructure/grpc/gen/health/v1"
	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/grpc/gen/health/v1/healthv1connect"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type HealthHandler struct {
	tracer trace.Tracer
}

var _ healthv1connect.HealthServiceHandler = (*HealthHandler)(nil)

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{tracer: otel.Tracer("connectrpc-server")}
}

func (h *HealthHandler) Check(
	ctx context.Context,
	req *connect.Request[v1.CheckRequest],
) (*connect.Response[v1.CheckResponse], error) {
	_, span := h.tracer.Start(ctx, "HealthHandler.Check")
	defer span.End()

	return connect.NewResponse(&v1.CheckResponse{
		Status: v1.ServingStatus_SERVING_STATUS_SERVING,
	}), nil
}
