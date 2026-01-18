package repository

import (
	"context"

	"github.com/Haya372/web-app-template/backend/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
}
