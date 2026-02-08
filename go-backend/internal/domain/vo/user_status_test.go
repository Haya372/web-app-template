package vo_test

import (
	"testing"

	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserStatusFromString(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect vo.UserStatus
	}{
		{name: "active", input: "ACTIVE", expect: vo.UserStatusActive},
		{name: "frozen", input: "FROZEN", expect: vo.UserStatusFrozen},
		{name: "deleted", input: "DELETED", expect: vo.UserStatusDeleted},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, err := vo.UserStatusFromString(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.expect, status)
		})
	}
}

func TestUserStatusFromString_Failure(t *testing.T) {
	_, err := vo.UserStatusFromString("UNKNOWN")

	require.Error(t, err)
}

func TestUserStatusPredicates(t *testing.T) {
	assert.True(t, vo.UserStatusActive.IsActive())
	assert.False(t, vo.UserStatusActive.IsFrozen())
	assert.False(t, vo.UserStatusActive.IsDeleted())

	assert.False(t, vo.UserStatusFrozen.IsActive())
	assert.True(t, vo.UserStatusFrozen.IsFrozen())
	assert.False(t, vo.UserStatusFrozen.IsDeleted())

	assert.False(t, vo.UserStatusDeleted.IsActive())
	assert.False(t, vo.UserStatusDeleted.IsFrozen())
	assert.True(t, vo.UserStatusDeleted.IsDeleted())
}
