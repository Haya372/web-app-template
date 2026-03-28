package http

import (
	"context"
	"errors"

	"github.com/Haya372/web-app-template/go-backend/internal/common"
	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	generated "github.com/Haya372/web-app-template/go-backend/internal/infrastructure/http/generated"
	commandpost "github.com/Haya372/web-app-template/go-backend/internal/usecase/command/post"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
)

// PostV1Posts handles POST /v1/posts (requires JWT).
func (h *serverHandler) PostV1Posts(
	ctx context.Context,
	req generated.PostV1PostsRequestObject,
) (generated.PostV1PostsResponseObject, error) {
	ctx, span := h.tracer.Start(ctx, "createPost")
	defer span.End()

	userIDStr := common.UserIDFromContext(ctx)
	if userIDStr == "" {
		h.logger.Error(ctx, "user ID missing from context — JWT middleware may not be applied")
		span.SetStatus(codes.Error, "missing user ID in context")

		return generated.PostV1Posts401ApplicationProblemPlusJSONResponse{
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

		return generated.PostV1Posts400ApplicationProblemPlusJSONResponse{
			BadRequestApplicationProblemPlusJSONResponse: generated.BadRequestApplicationProblemPlusJSONResponse(
				validationProblem("invalid user ID in token", nil),
			),
		}, nil
	}

	output, err := h.createPostUseCase.Execute(ctx, commandpost.CreatePostInput{
		UserID:  userID,
		Content: req.Body.Content,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return mapCreatePostError(err), nil
	}

	return generated.PostV1Posts201JSONResponse{
		Id:        output.ID,
		UserId:    output.UserID,
		Content:   output.Content,
		CreatedAt: output.CreatedAt,
	}, nil
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
		default:
		}
	}

	internalResp := generated.InternalServerErrorApplicationProblemPlusJSONResponse(internalProblem())

	return generated.PostV1Posts500ApplicationProblemPlusJSONResponse{
		InternalServerErrorApplicationProblemPlusJSONResponse: internalResp,
	}
}
