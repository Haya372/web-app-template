package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/Haya372/web-app-template/go-backend/internal/common"
	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	generated "github.com/Haya372/web-app-template/go-backend/internal/infrastructure/http/generated"
	commandpost "github.com/Haya372/web-app-template/go-backend/internal/usecase/command/post"
	commanduser "github.com/Haya372/web-app-template/go-backend/internal/usecase/command/user"
	queryuser "github.com/Haya372/web-app-template/go-backend/internal/usecase/query/user"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// serverHandler implements generated.StrictServerInterface and contains all
// HTTP handler logic for the API. It delegates business operations to use cases
// and maps domain errors to typed OpenAPI response objects.
type serverHandler struct {
	logger           common.Logger
	tracer           trace.Tracer
	signupUseCase    commanduser.SingupUseCase
	loginUseCase     commanduser.LoginUseCase
	listUsersUseCase queryuser.ListUsersUseCase
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

// PostV1UsersSignup handles POST /v1/users/signup.
func (h *serverHandler) PostV1UsersSignup(
	ctx context.Context,
	req generated.PostV1UsersSignupRequestObject,
) (generated.PostV1UsersSignupResponseObject, error) {
	ctx, span := h.tracer.Start(ctx, "signup")
	defer span.End()

	output, err := h.signupUseCase.Execute(ctx, commanduser.SignupInput{
		Email:    string(req.Body.Email),
		Password: req.Body.Password,
		Name:     req.Body.Name,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return mapSignupError(err), nil
	}

	return generated.PostV1UsersSignup201JSONResponse{
		Id:        openapi_types.UUID(output.Id),
		Name:      output.Name,
		Email:     openapi_types.Email(output.Email),
		Status:    output.Status.String(),
		CreatedAt: output.CreatedAt,
	}, nil
}

// PostV1UsersLogin handles POST /v1/users/login.
func (h *serverHandler) PostV1UsersLogin(
	ctx context.Context,
	req generated.PostV1UsersLoginRequestObject,
) (generated.PostV1UsersLoginResponseObject, error) {
	ctx, span := h.tracer.Start(ctx, "login")
	defer span.End()

	output, err := h.loginUseCase.Execute(ctx, commanduser.LoginInput{
		Email:    string(req.Body.Email),
		Password: req.Body.Password,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return mapLoginError(err), nil
	}

	return generated.PostV1UsersLogin200JSONResponse{
		Token:     output.Token,
		ExpiresAt: output.ExpiresAt,
		User: struct {
			Email openapi_types.Email `json:"email"`
			Id    string              `json:"id"`
			Name  string              `json:"name"`
		}{
			Id:    output.UserId,
			Name:  output.UserName,
			Email: openapi_types.Email(output.UserEmail),
		},
	}, nil
}

// GetV1Users handles GET /v1/users (requires JWT and users:list permission).
func (h *serverHandler) GetV1Users(
	ctx context.Context,
	req generated.GetV1UsersRequestObject,
) (generated.GetV1UsersResponseObject, error) {
	ctx, span := h.tracer.Start(ctx, "listUsers")
	defer span.End()

	userIdStr := common.UserIdFromContext(ctx)
	if userIdStr == "" {
		h.logger.Error(ctx, "user ID missing from context — JWT middleware may not be applied")
		span.SetStatus(codes.Error, "missing user ID in context")

		return generated.GetV1Users401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: generated.UnauthorizedApplicationProblemPlusJSONResponse(
				unauthorizedProblem(),
			),
		}, nil
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		parseErr := vo.NewValidationError("invalid user ID in token", nil, err)
		span.RecordError(parseErr)
		span.SetStatus(codes.Error, parseErr.Error())

		return generated.GetV1Users400ApplicationProblemPlusJSONResponse{
			BadRequestApplicationProblemPlusJSONResponse: generated.BadRequestApplicationProblemPlusJSONResponse(
				validationProblem("invalid user ID in token", nil),
			),
		}, nil
	}

	limit := 20
	if req.Params.Limit != nil {
		limit = *req.Params.Limit
	}

	offset := 0
	if req.Params.Offset != nil {
		offset = *req.Params.Offset
	}

	output, err := h.listUsersUseCase.Execute(ctx, queryuser.ListUsersInput{
		UserId: userId,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return mapListUsersError(err), nil
	}

	users := make([]generated.UserResponse, 0, len(output.Users))
	for _, u := range output.Users {
		users = append(users, generated.UserResponse{
			Id:        openapi_types.UUID(u.Id),
			Name:      u.Name,
			Email:     openapi_types.Email(u.Email),
			Status:    u.Status,
			CreatedAt: u.CreatedAt,
		})
	}

	return generated.GetV1Users200JSONResponse{
		Users:  users,
		Total:  output.Total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

// PostV1Posts handles POST /v1/posts (requires JWT).
func (h *serverHandler) PostV1Posts(
	ctx context.Context,
	req generated.PostV1PostsRequestObject,
) (generated.PostV1PostsResponseObject, error) {
	ctx, span := h.tracer.Start(ctx, "createPost")
	defer span.End()

	userIdStr := common.UserIdFromContext(ctx)
	if userIdStr == "" {
		h.logger.Error(ctx, "user ID missing from context — JWT middleware may not be applied")
		span.SetStatus(codes.Error, "missing user ID in context")

		return generated.PostV1Posts401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: generated.UnauthorizedApplicationProblemPlusJSONResponse(
				unauthorizedProblem(),
			),
		}, nil
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		parseErr := vo.NewValidationError("invalid user ID in token", nil, err)
		span.RecordError(parseErr)
		span.SetStatus(codes.Error, parseErr.Error())

		return generated.PostV1Posts400ApplicationProblemPlusJSONResponse{
			BadRequestApplicationProblemPlusJSONResponse: generated.BadRequestApplicationProblemPlusJSONResponse(
				validationProblem("invalid user ID in token", nil),
			),
		}, nil
	}

	output, err := h.createPostUseCase.Execute(ctx, commandpost.CreatePostInput{
		UserId:  userId,
		Content: req.Body.Content,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return mapCreatePostError(err), nil
	}

	return generated.PostV1Posts201JSONResponse{
		Id:        openapi_types.UUID(output.Id),
		UserId:    openapi_types.UUID(output.UserId),
		Content:   output.Content,
		CreatedAt: output.CreatedAt,
	}, nil
}

// --- error mapping helpers ---

func mapSignupError(err error) generated.PostV1UsersSignupResponseObject {
	var domainErr vo.Error
	if errors.As(err, &domainErr) {
		switch domainErr.Code() {
		case vo.ValidationErrorCode:
			return generated.PostV1UsersSignup400ApplicationProblemPlusJSONResponse{
				BadRequestApplicationProblemPlusJSONResponse: generated.BadRequestApplicationProblemPlusJSONResponse(
					validationProblemFromDomain(domainErr),
				),
			}
		case vo.DuplicateEmailErrorCode:
			return generated.PostV1UsersSignup409ApplicationProblemPlusJSONResponse{
				ConflictApplicationProblemPlusJSONResponse: generated.ConflictApplicationProblemPlusJSONResponse(
					domainErrToProblem(domainErr),
				),
			}
		}
	}

	return generated.PostV1UsersSignup500ApplicationProblemPlusJSONResponse{
		InternalServerErrorApplicationProblemPlusJSONResponse: generated.InternalServerErrorApplicationProblemPlusJSONResponse(
			internalProblem(),
		),
	}
}

func mapLoginError(err error) generated.PostV1UsersLoginResponseObject {
	var domainErr vo.Error
	if errors.As(err, &domainErr) {
		switch domainErr.Code() {
		case vo.ValidationErrorCode:
			return generated.PostV1UsersLogin400ApplicationProblemPlusJSONResponse{
				BadRequestApplicationProblemPlusJSONResponse: generated.BadRequestApplicationProblemPlusJSONResponse(
					validationProblemFromDomain(domainErr),
				),
			}
		case vo.InvalidCredentialErrorCode:
			return generated.PostV1UsersLogin401ApplicationProblemPlusJSONResponse{
				UnauthorizedApplicationProblemPlusJSONResponse: generated.UnauthorizedApplicationProblemPlusJSONResponse(
					domainErrToProblem(domainErr),
				),
			}
		}
	}

	return generated.PostV1UsersLogin500ApplicationProblemPlusJSONResponse{
		InternalServerErrorApplicationProblemPlusJSONResponse: generated.InternalServerErrorApplicationProblemPlusJSONResponse(
			internalProblem(),
		),
	}
}

func mapListUsersError(err error) generated.GetV1UsersResponseObject {
	var domainErr vo.Error
	if errors.As(err, &domainErr) {
		switch domainErr.Code() {
		case vo.ValidationErrorCode:
			return generated.GetV1Users400ApplicationProblemPlusJSONResponse{
				BadRequestApplicationProblemPlusJSONResponse: generated.BadRequestApplicationProblemPlusJSONResponse(
					validationProblemFromDomain(domainErr),
				),
			}
		case vo.UnauthorizedErrorCode, vo.InvalidCredentialErrorCode:
			return generated.GetV1Users401ApplicationProblemPlusJSONResponse{
				UnauthorizedApplicationProblemPlusJSONResponse: generated.UnauthorizedApplicationProblemPlusJSONResponse(
					domainErrToProblem(domainErr),
				),
			}
		case vo.ForbiddenErrorCode:
			return generated.GetV1Users403ApplicationProblemPlusJSONResponse{
				ForbiddenApplicationProblemPlusJSONResponse: generated.ForbiddenApplicationProblemPlusJSONResponse(
					domainErrToProblem(domainErr),
				),
			}
		}
	}

	return generated.GetV1Users500ApplicationProblemPlusJSONResponse{
		InternalServerErrorApplicationProblemPlusJSONResponse: generated.InternalServerErrorApplicationProblemPlusJSONResponse(
			internalProblem(),
		),
	}
}

func mapCreatePostError(err error) generated.PostV1PostsResponseObject {
	var domainErr vo.Error
	if errors.As(err, &domainErr) {
		switch domainErr.Code() {
		case vo.ValidationErrorCode:
			return generated.PostV1Posts400ApplicationProblemPlusJSONResponse{
				BadRequestApplicationProblemPlusJSONResponse: generated.BadRequestApplicationProblemPlusJSONResponse(
					validationProblemFromDomain(domainErr),
				),
			}
		case vo.UnauthorizedErrorCode, vo.InvalidCredentialErrorCode:
			return generated.PostV1Posts401ApplicationProblemPlusJSONResponse{
				UnauthorizedApplicationProblemPlusJSONResponse: generated.UnauthorizedApplicationProblemPlusJSONResponse(
					domainErrToProblem(domainErr),
				),
			}
		}
	}

	return generated.PostV1Posts500ApplicationProblemPlusJSONResponse{
		InternalServerErrorApplicationProblemPlusJSONResponse: generated.InternalServerErrorApplicationProblemPlusJSONResponse(
			internalProblem(),
		),
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

