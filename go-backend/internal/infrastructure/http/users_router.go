package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	"github.com/Haya372/web-app-template/go-backend/internal/usecase/command/user"
	queryuser "github.com/Haya372/web-app-template/go-backend/internal/usecase/query/user"
	"github.com/labstack/echo/v5"
	"go.opentelemetry.io/otel/codes"
)

const (
	defaultListLimit  = 20
	defaultListOffset = 0
)

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
			Id:    output.UserId,
			Name:  output.UserName,
			Email: output.UserEmail,
		},
	}

	return c.JSON(http.StatusOK, res)
}

func (r *routerImpl) handleListUsers(c *echo.Context) error {
	ctx := c.Request().Context()

	ctx, span := r.tracer.Start(ctx, "listUsers")
	defer span.End()

	limit := defaultListLimit
	offset := defaultListOffset

	if raw := c.QueryParam("limit"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			validationErr := vo.NewValidationError("limit must be an integer", nil, err)
			status, res := handleError(validationErr)
			span.RecordError(validationErr)
			span.SetStatus(codes.Error, validationErr.Error())

			return writeProblem(c, status, res)
		}

		limit = parsed
	}

	if raw := c.QueryParam("offset"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			validationErr := vo.NewValidationError("offset must be an integer", nil, err)
			status, res := handleError(validationErr)
			span.RecordError(validationErr)
			span.SetStatus(codes.Error, validationErr.Error())

			return writeProblem(c, status, res)
		}

		offset = parsed
	}

	output, err := r.ListUsersUseCase.Execute(ctx, queryuser.ListUsersInput{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		status, res := handleError(err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return writeProblem(c, status, res)
	}

	type userJSON struct {
		Id        string `json:"id"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		Status    string `json:"status"`
		CreatedAt string `json:"createdAt"`
	}

	users := make([]userJSON, 0, len(output.Users))
	for _, u := range output.Users {
		users = append(users, userJSON{
			Id:        u.Id.String(),
			Name:      u.Name,
			Email:     u.Email,
			Status:    u.Status,
			CreatedAt: u.CreatedAt.Format(time.RFC3339),
		})
	}

	return c.JSON(http.StatusOK, struct {
		Users  []userJSON `json:"users"`
		Total  int        `json:"total"`
		Limit  int        `json:"limit"`
		Offset int        `json:"offset"`
	}{
		Users:  users,
		Total:  output.Total,
		Limit:  limit,
		Offset: offset,
	})
}
