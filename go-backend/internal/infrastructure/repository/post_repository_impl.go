package repository

import (
	"context"

	"github.com/Haya372/web-app-template/go-backend/internal/common"
	"github.com/Haya372/web-app-template/go-backend/internal/domain/entity"
	"github.com/Haya372/web-app-template/go-backend/internal/domain/entity/repository"
	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/db"
	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/sqlc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type postRepositoryImpl struct {
	tracer    trace.Tracer
	logger    common.Logger
	dbManager db.DbManager
}

func (r *postRepositoryImpl) Create(ctx context.Context, post entity.Post) (entity.Post, error) {
	ctx, span := r.tracer.Start(ctx, "Create")
	defer span.End()

	var row sqlc.CreatePostRow

	err := r.dbManager.QueriesFunc(ctx, func(ctx context.Context, queries sqlc.Queries) error {
		var qErr error

		row, qErr = queries.CreatePost(ctx, sqlc.CreatePostParams{
			ID:        toPgtypeUuid(post.Id()),
			UserID:    toPgtypeUuid(post.UserId()),
			Content:   post.Content(),
			CreatedAt: toPgtypeTimestamp(post.CreatedAt()),
		})

		return qErr
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	return entity.ReconstructPost(
		row.ID.Bytes,
		row.UserID.Bytes,
		row.Content,
		row.CreatedAt.Time,
	), nil
}

func NewPostRepository(dbManager db.DbManager) repository.PostRepository {
	return &postRepositoryImpl{
		tracer:    otel.Tracer("PostRepository"),
		logger:    common.NewLogger(),
		dbManager: dbManager,
	}
}
