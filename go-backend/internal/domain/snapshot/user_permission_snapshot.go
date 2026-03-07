package snapshot

import (
	"slices"

	"github.com/Haya372/web-app-template/go-backend/internal/domain/entity"
	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	"github.com/google/uuid"
)

// UserPermissionSnapshot is a read-only projection combining a user's identity and their
// effective permissions, derived from all roles assigned to the user.
type UserPermissionSnapshot struct {
	UserId      uuid.UUID
	User        entity.User
	Permissions []vo.Permission
}

// HasPermission reports whether the snapshot grants the given permission.
func (s *UserPermissionSnapshot) HasPermission(p vo.Permission) bool {
	return slices.Contains(s.Permissions, p)
}
