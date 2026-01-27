package http

import (
	"context"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

const (
	gracefulTimeout = 5 * time.Second
)

type Server struct {
	Config echo.StartConfig
	Echo   *echo.Echo
}

func (s *Server) Start(ctx context.Context) error {
	return s.Config.Start(ctx, s.Echo)
}

func NewServer(r *Router) *echo.Echo {
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	r.AddRoute(e)

	return e
}

func NewEchoConfig() echo.StartConfig {
	return echo.StartConfig{
		Address:         ":8080",
		GracefulTimeout: gracefulTimeout,
	}
}
