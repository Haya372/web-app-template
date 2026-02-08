package http

import (
	"context"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const (
	gracefulTimeout = 5 * time.Second
)

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.ErrBadRequest.Wrap(err)
	}
	return nil
}

type Server struct {
	Config echo.StartConfig
	Echo   *echo.Echo
}

func (s *Server) Start(ctx context.Context) error {
	return s.Config.Start(ctx, s.Echo)
}

func NewServer(r *Router) *echo.Echo {
	e := echo.New()

	e.Validator = &customValidator{validator: validator.New()}

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	// TODO: replace otelecho middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			handler := otelhttp.NewHandler(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					c.SetRequest(r)
					_ = next(c)
				}),
				c.Path(),
			)

			handler.ServeHTTP(c.Response(), c.Request())

			return nil
		}
	})

	r.AddRoute(e)

	return e
}

func NewEchoConfig() echo.StartConfig {
	return echo.StartConfig{
		Address:         ":8080",
		GracefulTimeout: gracefulTimeout,
	}
}
