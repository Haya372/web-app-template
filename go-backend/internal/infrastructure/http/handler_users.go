package http

import (
	"context"
	"errors"

	"github.com/Haya372/web-app-template/go-backend/internal/common"
	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	generated "github.com/Haya372/web-app-template/go-backend/internal/infrastructure/http/generated"
	commanduser "github.com/Haya372/web-app-template/go-backend/internal/usecase/command/user"
	queryuser "github.com/Haya372/web-app-template/go-backend/internal/usecase/query/user"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"go.opentelemetry.io/otel/codes"
)

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
		Id:        output.ID,
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
			Id:    output.UserID,
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

	userIDStr := common.UserIDFromContext(ctx)
	if userIDStr == "" {
		h.logger.Error(ctx, "user ID missing from context — JWT middleware may not be applied")
		span.SetStatus(codes.Error, "missing user ID in context")

		return generated.GetV1Users401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: generated.UnauthorizedApplicationProblemPlusJSONResponse(
				unauthorizedProblem(),
			),
		}, nil
	}

	userID, err := uuid.Parse(userIDStr)
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
		UserID: userID,
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
			Id:        u.ID,
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
		default:
		}
	}

	internalResp := generated.InternalServerErrorApplicationProblemPlusJSONResponse(internalProblem())

	return generated.PostV1UsersSignup500ApplicationProblemPlusJSONResponse{
		InternalServerErrorApplicationProblemPlusJSONResponse: internalResp,
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
		default:
		}
	}

	internalResp := generated.InternalServerErrorApplicationProblemPlusJSONResponse(internalProblem())

	return generated.PostV1UsersLogin500ApplicationProblemPlusJSONResponse{
		InternalServerErrorApplicationProblemPlusJSONResponse: internalResp,
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
		default:
		}
	}

	internalResp := generated.InternalServerErrorApplicationProblemPlusJSONResponse(internalProblem())

	return generated.GetV1Users500ApplicationProblemPlusJSONResponse{
		InternalServerErrorApplicationProblemPlusJSONResponse: internalResp,
	}
}
