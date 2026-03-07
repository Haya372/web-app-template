package snapshot_test

import (
	"testing"

	"github.com/Haya372/web-app-template/go-backend/internal/domain/snapshot"
	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserPermissionSnapshot_HasPermission(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name        string
		permissions []vo.Permission
		check       vo.Permission
		expected    bool
	}{
		{
			name:        "returns true when permission is present",
			permissions: []vo.Permission{vo.PermissionUsersList},
			check:       vo.PermissionUsersList,
			expected:    true,
		},
		{
			name:        "returns false when permission is absent",
			permissions: []vo.Permission{},
			check:       vo.PermissionUsersList,
			expected:    false,
		},
		{
			name:        "returns false when snapshot has no permissions",
			permissions: nil,
			check:       vo.PermissionUsersList,
			expected:    false,
		},
		{
			name: "returns true when one of multiple permissions matches",
			permissions: []vo.Permission{
				vo.PermissionUsersCreate,
				vo.PermissionUsersList,
			},
			check:    vo.PermissionUsersList,
			expected: true,
		},
		{
			name: "returns false when none of multiple permissions match",
			permissions: []vo.Permission{
				vo.PermissionUsersCreate,
			},
			check:    vo.PermissionUsersList,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &snapshot.UserPermissionSnapshot{
				UserId:      userID,
				Permissions: tt.permissions,
			}
			assert.Equal(t, tt.expected, s.HasPermission(tt.check))
		})
	}
}
