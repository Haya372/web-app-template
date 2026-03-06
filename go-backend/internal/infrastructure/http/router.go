package http

import (
	"github.com/Haya372/web-app-template/go-backend/internal/common"
	"github.com/Haya372/web-app-template/go-backend/internal/usecase/command/user"
	queryuser "github.com/Haya372/web-app-template/go-backend/internal/usecase/query/user"
	"github.com/Haya372/web-app-template/go-backend/internal/usecase/service"
	"github.com/labstack/echo/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Router interface {
	AddRoute(e *echo.Echo)
}

type routerImpl struct {
	logger           common.Logger
	tracer           trace.Tracer
	SignupUseCase    user.SingupUseCase
	LoginUseCase     user.LoginUseCase
	ListUsersUseCase queryuser.ListUsersUseCase
	jwtService       service.JwtService
}

func (r *routerImpl) AddRoute(e *echo.Echo) {
	v1 := e.Group("/v1")
	v1.POST("/users/signup", r.handleSignup)
	v1.POST("/users/login", r.handleLogin)
	v1.GET("/users", r.handleListUsers, JWTMiddleware(r.jwtService))
}

func NewRouter(
	signupUseCase user.SingupUseCase,
	loginUseCase user.LoginUseCase,
	listUsersUseCase queryuser.ListUsersUseCase,
	jwtService service.JwtService,
) Router {
	return &routerImpl{
		logger:           common.NewLogger(),
		tracer:           otel.Tracer("root"),
		SignupUseCase:    signupUseCase,
		LoginUseCase:     loginUseCase,
		ListUsersUseCase: listUsersUseCase,
		jwtService:       jwtService,
	}
}

func writeProblem(c *echo.Context, status int, res problemDetails) error {
	c.Response().Header().Set(echo.HeaderContentType, problemContentType)

	return c.JSON(status, res)
}
