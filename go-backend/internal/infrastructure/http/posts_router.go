package http

import (
	"net/http"
	"time"

	"github.com/Haya372/web-app-template/go-backend/internal/common"
	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	"github.com/Haya372/web-app-template/go-backend/internal/usecase/command/post"
	"github.com/Haya372/web-app-template/go-backend/internal/usecase/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type postsRouter struct {
	logger              common.Logger
	tracer              trace.Tracer
	CreatePostUseCase   post.CreatePostUseCase
	jwtService          service.JwtService
}

func (r *postsRouter) AddRoute(e *echo.Echo) {
	g := e.Group("/v1/posts")
	g.POST("", r.handleCreatePost, JWTMiddleware(r.jwtService))
}

func (r *postsRouter) handleCreatePost(c *echo.Context) error {
	ctx := c.Request().Context()

	ctx, span := r.tracer.Start(ctx, "createPost")
	defer span.End()

	var req struct {
		Content string `json:"content" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		r.logger.Error(ctx, "failed to bind create post input", "error", err)
		status, res := handleError(err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return writeProblem(c, status, res)
	}

	if err := c.Validate(&req); err != nil {
		r.logger.Error(ctx, "failed to validate input", "error", err)
		status, res := handleError(err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return writeProblem(c, status, res)
	}

	userIdStr := common.UserIdFromContext(ctx)
	if userIdStr == "" {
		r.logger.Error(ctx, "user ID missing from context — JWT middleware may not be applied to this route")
		span.SetStatus(codes.Error, "missing user ID in context")

		return writeProblem(c, http.StatusUnauthorized, buildUnauthorizedProblem())
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		parseErr := vo.NewValidationError("invalid user ID in token", nil, err)
		status, res := handleError(parseErr)
		span.RecordError(parseErr)
		span.SetStatus(codes.Error, parseErr.Error())

		return writeProblem(c, status, res)
	}

	output, err := r.CreatePostUseCase.Execute(ctx, post.CreatePostInput{
		UserId:  userId,
		Content: req.Content,
	})
	if err != nil {
		status, res := handleError(err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return writeProblem(c, status, res)
	}

	res := struct {
		Id        string `json:"id"`
		UserId    string `json:"userId"`
		Content   string `json:"content"`
		CreatedAt string `json:"createdAt"`
	}{
		Id:        output.Id.String(),
		UserId:    output.UserId.String(),
		Content:   output.Content,
		CreatedAt: output.CreatedAt.Format(time.RFC3339),
	}

	return c.JSON(http.StatusCreated, res)
}
