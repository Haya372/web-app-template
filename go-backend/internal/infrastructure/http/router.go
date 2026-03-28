package http

import (
	stdhttp "net/http"

	generated "github.com/Haya372/web-app-template/go-backend/internal/infrastructure/http/generated"
	commandpost "github.com/Haya372/web-app-template/go-backend/internal/usecase/command/post"
	"github.com/Haya372/web-app-template/go-backend/internal/usecase/command/user"
	querypost "github.com/Haya372/web-app-template/go-backend/internal/usecase/query/post"
	queryuser "github.com/Haya372/web-app-template/go-backend/internal/usecase/query/user"
	"github.com/Haya372/web-app-template/go-backend/internal/usecase/service"
	"github.com/labstack/echo/v5"
)

// Router registers all HTTP routes on an Echo instance.
type Router interface {
	AddRoute(e *echo.Echo)
}

type routerImpl struct {
	handler    *serverHandler
	jwtService service.JwtService
}

func (r *routerImpl) AddRoute(e *echo.Echo) {
	// Health check — no authentication required.
	e.GET("/health", func(c *echo.Context) error {
		return c.NoContent(stdhttp.StatusOK)
	})

	// Build the strict handler (StrictServerInterface → ServerInterface adaptor).
	// The ServerInterfaceWrapper then adapts each ServerInterface method to the
	// standard net/http func(w, r) signature used below.
	//
	// oapi-codegen generates chi-server target code that uses net/http types, so
	// we adapt each method to Echo v5 here rather than using the chi router.
	// NewStrictHandlerWithOptions lets us replace the default RequestErrorHandlerFunc
	// (which uses http.Error / plain text) with our problem+json handler.
	strictHandler := generated.NewStrictHandlerWithOptions(r.handler, nil, generated.StrictHTTPServerOptions{
		RequestErrorHandlerFunc: apiErrorHandler,
		ResponseErrorHandlerFunc: func(w stdhttp.ResponseWriter, _ *stdhttp.Request, err error) {
			w.Header().Set("Content-Type", problemContentType)
			w.WriteHeader(stdhttp.StatusInternalServerError)
			_ = writeJSONError(w, err)
		},
	})

	siw := &generated.ServerInterfaceWrapper{
		Handler:          strictHandler,
		ErrorHandlerFunc: apiErrorHandler,
	}

	wrap := func(h func(stdhttp.ResponseWriter, *stdhttp.Request)) echo.HandlerFunc {
		return func(c *echo.Context) error {
			h(c.Response(), c.Request())

			return nil
		}
	}

	// Public routes.
	e.POST("/v1/users/signup", wrap(siw.PostV1UsersSignup))
	e.POST("/v1/users/login", wrap(siw.PostV1UsersLogin))

	// Protected routes — JWT validation is enforced by the middleware.
	e.GET("/v1/users", wrap(siw.GetV1Users), JWTMiddleware(r.jwtService))
	e.GET("/v1/posts", wrap(siw.GetV1Posts), JWTMiddleware(r.jwtService))
	e.POST("/v1/posts", wrap(siw.PostV1Posts), JWTMiddleware(r.jwtService))
}

// apiErrorHandler writes a problem+json error response for request-parse failures
// produced by the generated ServerInterfaceWrapper (e.g. invalid query param types).
func apiErrorHandler(w stdhttp.ResponseWriter, _ *stdhttp.Request, err error) {
	w.Header().Set("Content-Type", problemContentType)
	w.WriteHeader(stdhttp.StatusBadRequest)

	_ = writeJSONError(w, err)
}

// NewRouter constructs a Router backed by the generated StrictServerInterface.
func NewRouter(
	signupUseCase user.SingupUseCase,
	loginUseCase user.LoginUseCase,
	listUsersUseCase queryuser.ListUsersUseCase,
	createPostUseCase commandpost.CreatePostUseCase,
	listPostsUseCase querypost.ListPostsUseCase,
	jwtService service.JwtService,
) Router {
	return &routerImpl{
		handler: newServerHandler(
			signupUseCase,
			loginUseCase,
			listUsersUseCase,
			createPostUseCase,
			listPostsUseCase,
		),
		jwtService: jwtService,
	}
}
