package user

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

type listUsersUseCaseImpl struct {
	tracer           trace.Tracer
	logger           common.Logger
	userQueryService UserQueryService
}

func (uc *listUsersUseCaseImpl) Execute(ctx context.Context, input ListUsersInput) (*ListUsersOutput, error) {
	ctx, span := uc.tracer.Start(ctx, "execute")
	defer span.End()

	userID := common.UserIDFromContext(ctx)
	uc.logger.Info(ctx, "list users requested", "userID", userID, "limit", input.Limit, "offset", input.Offset)

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

	users, total, err := uc.userQueryService.FindAll(ctx, input.Limit, input.Offset)
	if err != nil {
		uc.logger.Error(ctx, "failed to find users", "error", err, "userID", userID)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	if users == nil {
		users = []UserDto{}
	}

	return &ListUsersOutput{
		Users: users,
		Total: total,
	}, nil
}

func NewListUsersUseCase(userQueryService UserQueryService) ListUsersUseCase {
	return &listUsersUseCaseImpl{
		tracer:           otel.Tracer("ListUsersUseCase"),
		logger:           common.NewLogger(),
		userQueryService: userQueryService,
	}
}
