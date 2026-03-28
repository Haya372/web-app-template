package aggregate

import (
	"slices"

	"github.com/Haya372/web-app-template/go-backend/internal/domain/entity"
	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	"github.com/google/uuid"
)

// UserPermissionAggregate is an aggregate combining a user's identity and their
// effective permissions, derived from all roles assigned to the user.
type UserPermissionAggregate struct {
	UserID      uuid.UUID
	User        entity.User
	Permissions []vo.Permission
}

// HasPermission reports whether the aggregate grants the given permission.
func (a *UserPermissionAggregate) HasPermission(p vo.Permission) bool {
	return slices.Contains(a.Permissions, p)
}
