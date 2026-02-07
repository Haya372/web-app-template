//go:generate mockgen -source=user_repository.go -destination=../../../../test/mock/domain/entity/repository/mock_user_repository.go

package repository

import (
	"context"

	"github.com/Haya372/web-app-template/go-backend/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user entity.User) (entity.User, error)
	FindByEmail(ctx context.Context, email string) (entity.User, error)
}
