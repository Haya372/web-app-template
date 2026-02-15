package http

import (
	"net/http"
	"time"

	"github.com/Haya372/web-app-template/go-backend/internal/common"
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
	LoginUseCase  user.LoginUseCase
}

func (r *routerImpl) AddRoute(e *echo.Echo) {
	v1 := e.Group("/v1")
	v1.POST("/users/signup", r.handleSignup)
	v1.POST("/users/login", r.handleLogin)
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

	output, err := r.SignupUseCase.Execute(ctx, user.SignupInput{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	})
	if err != nil {
		status, res := handleError(err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return writeProblem(c, status, res)
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

func (r *routerImpl) handleLogin(c *echo.Context) error {
	ctx := c.Request().Context()

	ctx, span := r.tracer.Start(ctx, "login")
	defer span.End()

	var req struct {
		Email    string `form:"email"    json:"email"    validate:"required,email"`
		Password string `form:"password" json:"password" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		r.logger.Error(ctx, "failed to bind login input", "error", err)
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

	output, err := r.LoginUseCase.Execute(ctx, user.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		status, res := handleError(err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return writeProblem(c, status, res)
	}

	res := struct {
		Token     string `json:"token"`
		ExpiresAt string `json:"expiresAt"`
		User      struct {
			Id    string `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"user"`
	}{
		Token:     output.Token,
		ExpiresAt: output.ExpiresAt.Format(time.RFC3339),
		User: struct {
			Id    string `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
		}{
			Id:    output.UserID,
			Name:  output.UserName,
			Email: output.UserEmail,
		},
	}

	return c.JSON(http.StatusOK, res)
}

func NewRouter(signupUseCase user.SingupUseCase, loginUseCase user.LoginUseCase) Router {
	return &routerImpl{
		logger:        common.NewLogger(),
		tracer:        otel.Tracer("root"),
		SignupUseCase: signupUseCase,
		LoginUseCase:  loginUseCase,
	}
}

func writeProblem(c *echo.Context, status int, res problemDetails) error {
	c.Response().Header().Set(echo.HeaderContentType, problemContentType)

	return c.JSON(status, res)
}
