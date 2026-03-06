package query

import (
	"context"
	"fmt"

	"github.com/Haya372/web-app-template/go-backend/internal/common"
	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/db"
	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/sqlc"
	usecasequery "github.com/Haya372/web-app-template/go-backend/internal/usecase/query/user"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type userQueryServiceImpl struct {
	tracer    trace.Tracer
	logger    common.Logger
	dbManager db.DbManager
}

func (s *userQueryServiceImpl) FindAll(ctx context.Context, limit, offset int) ([]usecasequery.UserDto, int, error) {
	ctx, span := s.tracer.Start(ctx, "FindAll")
	defer span.End()

	var (
		rows  []sqlc.FindAllUsersRow
		total int64
	)

	err := s.dbManager.QueriesFunc(ctx, func(ctx context.Context, queries sqlc.Queries) error {
		var err error

		rows, err = queries.FindAllUsers(ctx, sqlc.FindAllUsersParams{
			Limit:  int32(limit),  //nolint:gosec // limit is validated (1-100) by the use case layer
			Offset: int32(offset), //nolint:gosec // offset is validated (>=0) by the use case layer
		})
		if err != nil {
			return err
		}

		total, err = queries.CountUsers(ctx)

		return err
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		s.logger.Error(ctx, "failed to query users", "error", err)

		return nil, 0, err
	}

	dtos := make([]usecasequery.UserDto, 0, len(rows))

	for _, row := range rows {
		status, err := vo.UserStatusFromString(row.StatusCode)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())

			return nil, 0, fmt.Errorf("parse user status: %w", err)
		}

		dtos = append(dtos, usecasequery.UserDto{
			Id:        uuid.UUID(row.ID.Bytes),
			Name:      row.Name,
			Email:     row.Email,
			Status:    status.String(),
			CreatedAt: row.CreatedAt.Time,
		})
	}

	return dtos, int(total), nil
}

func NewUserQueryService(dbManager db.DbManager) usecasequery.UserQueryService {
	return &userQueryServiceImpl{
		tracer:    otel.Tracer("UserQueryService"),
		logger:    common.NewLogger(),
		dbManager: dbManager,
	}
}
