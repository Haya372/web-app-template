package user

import (
	"context"
	"errors"

	"github.com/Haya372/web-app-template/go-backend/internal/common"
	"github.com/Haya372/web-app-template/go-backend/internal/domain/snapshot/reader"
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
	errInvalidLimit      = errors.New("limit out of range")
	errInvalidOffset     = errors.New("offset must be non-negative")
	errLacksUsersListPerm = errors.New("user lacks users:list permission")
)

type listUsersUseCaseImpl struct {
	tracer           trace.Tracer
	logger           common.Logger
	userQueryService UserQueryService
	permissionReader reader.UserPermissionReader
}

func (uc *listUsersUseCaseImpl) Execute(
	ctx context.Context, input ListUsersInput,
) (*ListUsersOutput, error) {
	ctx, span := uc.tracer.Start(ctx, "execute")
	defer span.End()

	uc.logger.Info(ctx, "list users requested", "limit", input.Limit, "offset", input.Offset)

	snap, err := uc.permissionReader.FindByUserId(ctx, input.UserId)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	if !snap.HasPermission(vo.PermissionUsersList) {
		err = vo.NewForbiddenError("insufficient permissions", nil, errLacksUsersListPerm)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

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
		uc.logger.Error(ctx, "failed to find users", "error", err)
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

func NewListUsersUseCase(
	userQueryService UserQueryService, permissionReader reader.UserPermissionReader,
) ListUsersUseCase {
	return &listUsersUseCaseImpl{
		tracer:           otel.Tracer("ListUsersUseCase"),
		logger:           common.NewLogger(),
		userQueryService: userQueryService,
		permissionReader: permissionReader,
	}
}
