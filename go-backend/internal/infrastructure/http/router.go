package http

import (
	"github.com/Haya372/web-app-template/go-backend/internal/common"
	"github.com/Haya372/web-app-template/go-backend/internal/usecase/command/user"
	queryuser "github.com/Haya372/web-app-template/go-backend/internal/usecase/query/user"
	"github.com/Haya372/web-app-template/go-backend/internal/usecase/service"
	"github.com/labstack/echo/v5"
	"go.opentelemetry.io/otel"
)

type Router interface {
	AddRoute(e *echo.Echo)
}

type routerImpl struct {
	usersRouter *usersRouter
}

func (r *routerImpl) AddRoute(e *echo.Echo) {
	r.usersRouter.AddRoute(e)
}

func NewRouter(
	signupUseCase user.SingupUseCase,
	loginUseCase user.LoginUseCase,
	listUsersUseCase queryuser.ListUsersUseCase,
	jwtService service.JwtService,
) Router {
	return &routerImpl{
		usersRouter: &usersRouter{
			logger:           common.NewLogger(),
			tracer:           otel.Tracer("users"),
			SignupUseCase:    signupUseCase,
			LoginUseCase:     loginUseCase,
			ListUsersUseCase: listUsersUseCase,
			jwtService:       jwtService,
		},
	}
}

func writeProblem(c *echo.Context, status int, res problemDetails) error {
	c.Response().Header().Set(echo.HeaderContentType, problemContentType)

	return c.JSON(status, res)
}
