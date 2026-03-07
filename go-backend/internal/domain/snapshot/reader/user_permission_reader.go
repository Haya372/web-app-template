//go:generate mockgen -source=user_permission_reader.go -destination=../../../../test/mock/domain/snapshot/reader/mock_user_permission_reader.go

package reader

import (
	"context"

	"github.com/Haya372/web-app-template/go-backend/internal/domain/snapshot"
	"github.com/google/uuid"
)

// UserPermissionReader is the port for fetching a user's effective permission snapshot.
type UserPermissionReader interface {
	FindByUserId(ctx context.Context, userId uuid.UUID) (*snapshot.UserPermissionSnapshot, error)
}
