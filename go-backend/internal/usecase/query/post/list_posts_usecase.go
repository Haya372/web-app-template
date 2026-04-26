package post

import (
	"context"
	"errors"

	"github.com/Haya372/web-app-template/go-backend/internal/common"
	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	minLimit  = 1
	maxLimit  = 100
	minOffset = 0
)

var (
	errInvalidLimit  = errors.New("limit out of range")
	errInvalidOffset = errors.New("offset must be non-negative")
)

type listPostsUseCaseImpl struct {
	tracer           trace.Tracer
	logger           common.Logger
	postQueryService PostQueryService
}

func (uc *listPostsUseCaseImpl) Execute(
	ctx context.Context, input ListPostsInput,
) (*ListPostsOutput, error) {
	ctx, span := uc.tracer.Start(ctx, "list_posts")
	defer span.End()

	uc.logger.Info(ctx, "list posts requested", "limit", input.Limit, "offset", input.Offset)

	if input.Limit < minLimit || input.Limit > maxLimit {
		err := vo.NewValidationError("limit must be between 1 and 100", nil, errInvalidLimit)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	if input.Offset < minOffset {
		err := vo.NewValidationError("offset must be 0 or greater", nil, errInvalidOffset)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	posts, total, err := uc.postQueryService.FindAll(ctx, input.Limit, input.Offset)
	if err != nil {
		uc.logger.Error(ctx, "failed to find posts", "error", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	if posts == nil {
		posts = []PostDto{}
	}

	return &ListPostsOutput{
		Posts: posts,
		Total: total,
	}, nil
}

// NewListPostsUseCase creates a new ListPostsUseCase.
func NewListPostsUseCase(postQueryService PostQueryService) ListPostsUseCase {
	return &listPostsUseCaseImpl{
		tracer:           otel.Tracer("ListPostsUseCase"),
		logger:           common.NewLogger(),
		postQueryService: postQueryService,
	}
}
