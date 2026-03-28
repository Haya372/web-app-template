package query

import (
	"context"
	"errors"
	"fmt"

	"github.com/Haya372/web-app-template/go-backend/internal/common"
	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/db"
	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/sqlc"
	usecasequery "github.com/Haya372/web-app-template/go-backend/internal/usecase/query/post"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var (
	errNullPostUUID      = errors.New("post row contains NULL uuid")
	errNullPostCreatedAt = errors.New("post row contains NULL created_at")
)

type postQueryServiceImpl struct {
	tracer    trace.Tracer
	logger    common.Logger
	dbManager db.DbManager
}

func (s *postQueryServiceImpl) FindAll(ctx context.Context, limit, offset int) ([]usecasequery.PostDto, int, error) {
	ctx, span := s.tracer.Start(ctx, "FindAll")
	defer span.End()

	var (
		rows  []sqlc.FindAllPostsRow
		total int64
	)

	err := s.dbManager.QueriesFunc(ctx, func(ctx context.Context, queries sqlc.Queries) error {
		var err error

		rows, err = queries.FindAllPosts(ctx, sqlc.FindAllPostsParams{
			Limit:  int32(limit),  //nolint:gosec // limit is validated (1-100) by the use case layer
			Offset: int32(offset), //nolint:gosec // offset is validated (>=0) by the use case layer
		})
		if err != nil {
			return err
		}

		total, err = queries.CountPosts(ctx)

		return err
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		s.logger.Error(ctx, "failed to query posts", "error", err)

		return nil, 0, err
	}

	dtos := make([]usecasequery.PostDto, 0, len(rows))

	for _, row := range rows {
		if !row.ID.Valid || !row.UserID.Valid {
			err := fmt.Errorf("%w: id_valid=%v user_id_valid=%v", errNullPostUUID, row.ID.Valid, row.UserID.Valid)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())

			return nil, 0, err
		}

		if !row.CreatedAt.Valid {
			err := fmt.Errorf("%w for id %s", errNullPostCreatedAt, uuid.UUID(row.ID.Bytes))
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())

			return nil, 0, err
		}

		dtos = append(dtos, usecasequery.PostDto{
			ID:        uuid.UUID(row.ID.Bytes),
			UserID:    uuid.UUID(row.UserID.Bytes),
			Content:   row.Content,
			CreatedAt: row.CreatedAt.Time,
		})
	}

	return dtos, int(total), nil
}

// NewPostQueryService creates a new PostQueryService backed by Postgres.
func NewPostQueryService(dbManager db.DbManager) usecasequery.PostQueryService {
	return &postQueryServiceImpl{
		tracer:    otel.Tracer("PostQueryService"),
		logger:    common.NewLogger(),
		dbManager: dbManager,
	}
}
