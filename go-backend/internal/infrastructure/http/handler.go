package http

import (
	"net/http"

	"github.com/Haya372/web-app-template/go-backend/internal/common"
	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	generated "github.com/Haya372/web-app-template/go-backend/internal/infrastructure/http/generated"
	commandpost "github.com/Haya372/web-app-template/go-backend/internal/usecase/command/post"
	commanduser "github.com/Haya372/web-app-template/go-backend/internal/usecase/command/user"
	queryuser "github.com/Haya372/web-app-template/go-backend/internal/usecase/query/user"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// serverHandler implements generated.StrictServerInterface and contains all
// HTTP handler logic for the API. It delegates business operations to use cases
// and maps domain errors to typed OpenAPI response objects.
type serverHandler struct {
	logger            common.Logger
	tracer            trace.Tracer
	signupUseCase     commanduser.SingupUseCase
	loginUseCase      commanduser.LoginUseCase
	listUsersUseCase  queryuser.ListUsersUseCase
	createPostUseCase commandpost.CreatePostUseCase
}

// Compile-time assertion that serverHandler satisfies the generated interface.
var _ generated.StrictServerInterface = (*serverHandler)(nil)

func newServerHandler(
	signupUseCase commanduser.SingupUseCase,
	loginUseCase commanduser.LoginUseCase,
	listUsersUseCase queryuser.ListUsersUseCase,
	createPostUseCase commandpost.CreatePostUseCase,
) *serverHandler {
	return &serverHandler{
		logger:            common.NewLogger(),
		tracer:            otel.Tracer("server"),
		signupUseCase:     signupUseCase,
		loginUseCase:      loginUseCase,
		listUsersUseCase:  listUsersUseCase,
		createPostUseCase: createPostUseCase,
	}
}

// --- problem detail builders ---

func domainErrToProblem(e vo.Error) generated.ProblemDetails {
	p := generated.ProblemDetails{
		Type:   string(e.Code()),
		Title:  e.Code().Title(),
		Status: e.Status(),
	}

	if msg := e.Message(); msg != "" {
		p.Detail = &msg
	}

	if details := e.Details(); len(details) > 0 {
		m := make(map[string][]string, len(details))
		for k, v := range details {
			if s, ok := v.(string); ok {
				m[k] = []string{s}
			}
		}

		if len(m) > 0 {
			p.Errors = &m
		}
	}

	return p
}

func validationProblemFromDomain(e vo.Error) generated.ProblemDetails {
	p := generated.ProblemDetails{
		Type:   string(vo.ValidationErrorCode),
		Title:  vo.ValidationErrorCode.Title(),
		Status: http.StatusBadRequest,
	}

	detail := "invalid request parameters"
	p.Detail = &detail

	if details := e.Details(); len(details) > 0 {
		m := make(map[string][]string, len(details))
		for k, v := range details {
			if s, ok := v.(string); ok {
				m[k] = []string{s}
			}
		}

		if len(m) > 0 {
			p.Errors = &m
		}
	}

	return p
}

func validationProblem(detail string, errors map[string][]string) generated.ProblemDetails {
	p := generated.ProblemDetails{
		Type:   string(vo.ValidationErrorCode),
		Title:  vo.ValidationErrorCode.Title(),
		Status: http.StatusBadRequest,
		Detail: &detail,
	}

	if len(errors) > 0 {
		p.Errors = &errors
	}

	return p
}

func unauthorizedProblem() generated.ProblemDetails {
	title := vo.UnauthorizedErrorCode.Title()

	return generated.ProblemDetails{
		Type:   string(vo.UnauthorizedErrorCode),
		Title:  title,
		Status: http.StatusUnauthorized,
	}
}

func internalProblem() generated.ProblemDetails {
	return generated.ProblemDetails{
		Type:   string(vo.InternalErrorCode),
		Title:  vo.InternalErrorCode.Title(),
		Status: http.StatusInternalServerError,
	}
}
