package repository

import (
	"context"

	"github.com/Haya372/web-app-template/backend/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user entity.User) (entity.User, error)
	FindByEmail(ctx context.Context, email string) (entity.User, error)
}
