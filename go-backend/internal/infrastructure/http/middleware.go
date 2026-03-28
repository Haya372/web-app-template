package http

import (
	"net/http"
	"strings"

	"github.com/Haya372/web-app-template/go-backend/internal/common"
	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	generated "github.com/Haya372/web-app-template/go-backend/internal/infrastructure/http/generated"
	"github.com/Haya372/web-app-template/go-backend/internal/usecase/service"
	"github.com/labstack/echo/v5"
)

// JWTMiddleware returns an Echo middleware that validates Bearer JWT tokens.
// On success the authenticated user's ID is stored in both the Echo context
// (key "userID") and the Go request context via common.WithUserId, so that
// downstream handlers and use cases can retrieve it.
// Requests without a valid token receive a 401 Unauthorized problem response.
func JWTMiddleware(jwtService service.JwtService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return writeUnauthorized(c)
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := jwtService.ValidateToken(c.Request().Context(), token)
			if err != nil {
				return writeUnauthorized(c)
			}

			// Propagate userID into both the Echo context and the Go request
			// context so it is available to use cases for logging.
			c.Set("userID", claims.UserID)
			ctx := common.WithUserID(c.Request().Context(), claims.UserID)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

func writeUnauthorized(c *echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, problemContentType)

	return c.JSON(http.StatusUnauthorized, generated.ProblemDetails{
		Type:   string(vo.UnauthorizedErrorCode),
		Title:  vo.UnauthorizedErrorCode.Title(),
		Status: http.StatusUnauthorized,
	})
}
