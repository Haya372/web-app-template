package reader

import (
	"context"
	"fmt"

	"github.com/Haya372/web-app-template/go-backend/internal/common"
	"github.com/Haya372/web-app-template/go-backend/internal/domain/entity"
	"github.com/Haya372/web-app-template/go-backend/internal/domain/snapshot"
	snapshotreader "github.com/Haya372/web-app-template/go-backend/internal/domain/snapshot/reader"
	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/db"
	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.opentelemetry.io/otel"
	otelcodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type userPermissionReaderImpl struct {
	tracer    trace.Tracer
	logger    common.Logger
	dbManager db.DbManager
}

func (r *userPermissionReaderImpl) FindByUserId(
	ctx context.Context, userId uuid.UUID,
) (*snapshot.UserPermissionSnapshot, error) {
	ctx, span := r.tracer.Start(ctx, "FindByUserId")
	defer span.End()

	pgID := pgtype.UUID{Bytes: userId, Valid: true}

	var (
		userRow   sqlc.User
		permCodes []string
	)

	err := r.dbManager.QueriesFunc(ctx, func(ctx context.Context, queries sqlc.Queries) error {
		var err error

		userRow, err = queries.FindUserByID(ctx, pgID)
		if err != nil {
			return fmt.Errorf("find user: %w", err)
		}

		permCodes, err = queries.FindPermissionsByUserID(ctx, pgID)

		return err
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelcodes.Error, err.Error())
		r.logger.Error(ctx, "failed to query user permission snapshot", "error", err)

		return nil, err
	}

	status, err := vo.UserStatusFromString(userRow.StatusCode)
	if err != nil {
		return nil, fmt.Errorf("parse user status: %w", err)
	}

	user := entity.ReconstructUser(
		uuid.UUID(userRow.ID.Bytes),
		userRow.Email,
		userRow.PasswordHash,
		userRow.Name,
		status,
		userRow.CreatedAt.Time,
	)

	perms := make([]vo.Permission, 0, len(permCodes))

	for _, code := range permCodes {
		perms = append(perms, vo.Permission(code))
	}

	return &snapshot.UserPermissionSnapshot{
		UserId:      userId,
		User:        user,
		Permissions: perms,
	}, nil
}

func NewUserPermissionReader(dbManager db.DbManager) snapshotreader.UserPermissionReader {
	return &userPermissionReaderImpl{
		tracer:    otel.Tracer("UserPermissionReader"),
		logger:    common.NewLogger(),
		dbManager: dbManager,
	}
}
