package http

import (
	"context"
	"errors"

	"github.com/Haya372/web-app-template/go-backend/internal/common"
	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	generated "github.com/Haya372/web-app-template/go-backend/internal/infrastructure/http/generated"
	commandpost "github.com/Haya372/web-app-template/go-backend/internal/usecase/command/post"
	querypost "github.com/Haya372/web-app-template/go-backend/internal/usecase/query/post"
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

// GetV1Posts handles GET /v1/posts (requires JWT).
func (h *serverHandler) GetV1Posts(
	ctx context.Context,
	req generated.GetV1PostsRequestObject,
) (generated.GetV1PostsResponseObject, error) {
	ctx, span := h.tracer.Start(ctx, "listPosts")
	defer span.End()

	if common.UserIDFromContext(ctx) == "" {
		h.logger.Error(ctx, "user ID missing from context — JWT middleware may not be applied")
		span.SetStatus(codes.Error, "missing user ID in context")

		return generated.GetV1Posts401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: generated.UnauthorizedApplicationProblemPlusJSONResponse(
				unauthorizedProblem(),
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

	output, err := h.listPostsUseCase.Execute(ctx, querypost.ListPostsInput{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return mapListPostsError(err), nil
	}

	posts := make([]generated.PostResponse, len(output.Posts))
	for i, p := range output.Posts {
		posts[i] = generated.PostResponse{
			Id:        p.ID,
			UserId:    p.UserID,
			Content:   p.Content,
			CreatedAt: p.CreatedAt,
		}
	}

	return generated.GetV1Posts200JSONResponse{
		Posts:  posts,
		Total:  output.Total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func mapListPostsError(err error) generated.GetV1PostsResponseObject {
	var domainErr vo.Error
	if errors.As(err, &domainErr) {
		switch domainErr.Code() {
		case vo.ValidationErrorCode:
			return generated.GetV1Posts400ApplicationProblemPlusJSONResponse{
				BadRequestApplicationProblemPlusJSONResponse: generated.BadRequestApplicationProblemPlusJSONResponse(
					validationProblemFromDomain(domainErr),
				),
			}
		default:
			// No 401/403 arm — this endpoint performs no authorization checks per Issue #58.
			// JWT auth is handled by JWTMiddleware before the handler is invoked.
		}
	}

	internalResp := generated.InternalServerErrorApplicationProblemPlusJSONResponse(internalProblem())

	return generated.GetV1Posts500ApplicationProblemPlusJSONResponse{
		InternalServerErrorApplicationProblemPlusJSONResponse: internalResp,
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
		default:
		}
	}

	internalResp := generated.InternalServerErrorApplicationProblemPlusJSONResponse(internalProblem())

	return generated.PostV1Posts500ApplicationProblemPlusJSONResponse{
		InternalServerErrorApplicationProblemPlusJSONResponse: internalResp,
	}
}
