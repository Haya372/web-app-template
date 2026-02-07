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
)

type Router struct {
	SignupUseCase user.SingupUseCase
}

func (r *Router) AddRoute(e *echo.Echo) {
	logger := common.NewLogger()
	tracer := otel.Tracer("root")

	e.POST("/signup", func(c *echo.Context) error {
		ctx := c.Request().Context()

		ctx, span := tracer.Start(ctx, "signup")
		defer span.End()

		var req struct {
			Email    string `form:"email"    json:"email"`
			Password string `form:"password" json:"password"`
			Name     string `form:"name"     json:"name"`
		}

		if err := c.Bind(&req); err != nil {
			logger.Error(ctx, "failed to bind signup input", "error", err)

			res := map[string]string{
				"code": string(vo.ValidationErrorCode),
			}

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
			CreatedAt string `json:"createdAt"`
		}{
			Id:        output.Id.URN(),
			Name:      output.Name,
			Email:     output.Email,
			CreatedAt: output.CreatedAt.Format(time.RFC3339),
		}

		return c.JSON(http.StatusCreated, res)
	})
}
