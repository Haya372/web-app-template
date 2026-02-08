package http

import (
	"net/http"
	"time"

	"github.com/Haya372/web-app-template/go-backend/internal/common"
	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	"github.com/Haya372/web-app-template/go-backend/internal/usecase/command/user"
	"github.com/labstack/echo/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Router interface {
	AddRoute(e *echo.Echo)
}

type routerImpl struct {
	logger        common.Logger
	tracer        trace.Tracer
	SignupUseCase user.SingupUseCase
}

func (r *routerImpl) AddRoute(e *echo.Echo) {
	e.POST("/signup", r.handleSignup)
}

func (r *routerImpl) handleSignup(c *echo.Context) error {
	ctx := c.Request().Context()

	ctx, span := r.tracer.Start(ctx, "signup")
	defer span.End()

	var req struct {
		Email    string `form:"email"    json:"email"    validate:"required,email"`
		Password string `form:"password" json:"password" validate:"required"`
		Name     string `form:"name"     json:"name"     validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		r.logger.Error(ctx, "failed to bind signup input", "error", err)

		res := map[string]string{
			"code": string(vo.ValidationErrorCode),
		}

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return c.JSON(http.StatusBadRequest, res)
	}

	if err := c.Validate(&req); err != nil {
		r.logger.Error(ctx, "failed to validate input", "error", err)

		res := map[string]string{
			"code": string(vo.ValidationErrorCode),
		}

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return c.JSON(http.StatusBadRequest, res)
	}

	output, err := r.SignupUseCase.Execute(ctx, user.SignupInput{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	})
	if err != nil {
		status, res := handleError(err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return c.JSON(status, res)
	}

	res := struct {
		Id        string `json:"id"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		Status    string `json:"status"`
		CreatedAt string `json:"createdAt"`
	}{
		Id:        output.Id.URN(),
		Name:      output.Name,
		Email:     output.Email,
		Status:    output.Status.String(),
		CreatedAt: output.CreatedAt.Format(time.RFC3339),
	}

	return c.JSON(http.StatusCreated, res)
}

func NewRouter(signupUseCase user.SingupUseCase) Router {
	return &routerImpl{
		logger:        common.NewLogger(),
		tracer:        otel.Tracer("root"),
		SignupUseCase: signupUseCase,
	}
}
