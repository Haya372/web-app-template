package http

import (
	"net/http"
	"time"

	"github.com/Haya372/web-app-template/backend/internal/common"
	"github.com/Haya372/web-app-template/backend/internal/domain/vo"
	"github.com/Haya372/web-app-template/backend/internal/usecase/command/user"
	"github.com/labstack/echo/v5"
)

type Router struct {
	SignupUseCase user.SingupUseCase
}

func (r *Router) AddRoute(e *echo.Echo) {
	logger := common.NewLogger()

	e.POST("/signup", func(c *echo.Context) error {
		ctx := c.Request().Context()

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
