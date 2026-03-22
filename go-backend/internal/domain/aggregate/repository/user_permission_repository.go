//go:generate mockgen -source=user_permission_repository.go -destination=../../../../test/mock/domain/aggregate/repository/mock_user_permission_repository.go

package repository

import (
	"context"

	"github.com/Haya372/web-app-template/go-backend/internal/domain/aggregate"
	"github.com/google/uuid"
)

// UserPermissionRepository is the port for fetching a user's effective permission aggregate.
type UserPermissionRepository interface {
	FindByUserId(ctx context.Context, userId uuid.UUID) (*aggregate.UserPermissionAggregate, error)
}
