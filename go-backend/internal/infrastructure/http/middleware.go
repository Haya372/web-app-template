package http

import (
	"net/http"
	"strings"

	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	"github.com/Haya372/web-app-template/go-backend/internal/usecase/service"
	"github.com/labstack/echo/v5"
)

// contextKey is an unexported type for context keys in this package.
type contextKey string

const userIDKey contextKey = "userID"

// JWTMiddleware returns an Echo middleware that validates Bearer JWT tokens.
// On success the authenticated user's ID is stored in the Echo context under
// the key "userID" so that downstream handlers can retrieve it via
// c.Get("userID").(string).
// Requests without a valid token receive a 401 Unauthorized problem response.
func JWTMiddleware(jwtService service.JwtService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return writeProblem(c, http.StatusUnauthorized, buildUnauthorizedProblem())
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := jwtService.ValidateToken(c.Request().Context(), token)
			if err != nil {
				return writeProblem(c, http.StatusUnauthorized, buildUnauthorizedProblem())
			}

			c.Set(string(userIDKey), claims.UserID)

			return next(c)
		}
	}
}

func buildUnauthorizedProblem() problemDetails {
	return problemDetails{
		Type:   vo.UnauthorizedErrorCode,
		Title:  vo.UnauthorizedErrorCode.Title(),
		Status: http.StatusUnauthorized,
	}
}
