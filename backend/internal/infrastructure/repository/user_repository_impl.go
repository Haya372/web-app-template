package repository

import (
	"context"

	"github.com/Haya372/go-template/backend/internal/common"
	"github.com/Haya372/go-template/backend/internal/domain"
	"github.com/Haya372/go-template/backend/internal/domain/repository"
	"github.com/Haya372/go-template/backend/internal/infrastructure/db"
	"github.com/Haya372/go-template/backend/internal/infrastructure/sqlc"
)

type userRepositoryImpl struct {
	logger    common.Logger
	dbManager db.DbManager
}

func (r *userRepositoryImpl) Create(ctx context.Context, user domain.User) (domain.User, error) {
	err := r.dbManager.QueriesFunc(ctx, func(ctx context.Context, queries sqlc.Queries) error {
		return queries.CreateUser(ctx, sqlc.CreateUserParams{
			ID:           toPgtypeUuid(user.Id()),
			Email:        user.Email(),
			PasswordHash: user.PasswordHash(),
			Name:         user.Name(),
			CreatedAt:    toPgtypeTimestamp(user.CreatedAt()),
		})
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	var user *sqlc.User
	err := r.dbManager.QueriesFunc(ctx, func(ctx context.Context, queries sqlc.Queries) error {
		u, err := queries.FindUserByEmail(ctx, email)
		if err != nil {
			return err
		}
		user = &u
		return nil
	})
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, err
	}

	return domain.ReconstructUser(
		user.ID.Bytes,
		user.Email,
		user.PasswordHash,
		user.Name,
		user.CreatedAt.Time,
	), nil
}

func NewUserRepository(dbManager db.DbManager) repository.UserRepository {
	return &userRepositoryImpl{
		logger:    common.NewLogger(),
		dbManager: dbManager,
	}
}
