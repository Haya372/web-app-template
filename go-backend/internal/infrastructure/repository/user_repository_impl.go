package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Haya372/web-app-template/go-backend/internal/common"
	"github.com/Haya372/web-app-template/go-backend/internal/domain/entity"
	"github.com/Haya372/web-app-template/go-backend/internal/domain/entity/repository"
	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/db"
	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/sqlc"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type userRepositoryImpl struct {
	tracer    trace.Tracer
	logger    common.Logger
	dbManager db.DbManager
}

func (r *userRepositoryImpl) Create(ctx context.Context, user entity.User) (entity.User, error) {
	ctx, span := r.tracer.Start(ctx, "Create")
	defer span.End()

	err := r.dbManager.QueriesFunc(ctx, func(ctx context.Context, queries sqlc.Queries) error {
		return queries.CreateUser(ctx, sqlc.CreateUserParams{
			ID:           toPgtypeUuid(user.Id()),
			Email:        user.Email(),
			PasswordHash: user.PasswordHash(),
			Name:         user.Name(),
			StatusCode:   user.Status().String(),
			CreatedAt:    toPgtypeTimestamp(user.CreatedAt()),
		})
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	return user, nil
}

func (r *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (entity.User, error) {
	ctx, span := r.tracer.Start(ctx, "FindByEmail")
	defer span.End()

	var dbUser *sqlc.User

	err := r.dbManager.QueriesFunc(ctx, func(ctx context.Context, queries sqlc.Queries) error {
		u, err := queries.FindUserByEmail(ctx, email)
		if err != nil {
			return err
		}

		dbUser = &u

		return nil
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrUserNotFound
		}

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	// NOTE: schema guarantees FK integrity, but keep this guard defensive to surface unexpected records.
	status, err := vo.UserStatusFromString(dbUser.StatusCode)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, fmt.Errorf("parse user status: %w", err)
	}

	return entity.ReconstructUser(
		dbUser.ID.Bytes,
		dbUser.Email,
		dbUser.PasswordHash,
		dbUser.Name,
		status,
		dbUser.CreatedAt.Time,
	), nil
}

func NewUserRepository(dbManager db.DbManager) repository.UserRepository {
	return &userRepositoryImpl{
		tracer:    otel.Tracer("UserRepository"),
		logger:    common.NewLogger(),
		dbManager: dbManager,
	}
}
