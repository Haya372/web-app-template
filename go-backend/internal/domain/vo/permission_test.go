package vo_test

import (
	"strings"
	"testing"

	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPermission_HappyCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "valid permission",
			input: "users:list",
		},
		{
			name:  "valid permission with create action",
			input: "users:create",
		},
		{
			name:  "boundary: exactly 128-char permission",
			input: strings.Repeat("a", 63) + ":" + strings.Repeat("b", 64), // 63 + 1 + 64 = 128
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			permission, err := vo.NewPermission(tt.input)

			require.NoError(t, err)
			assert.NotNil(t, permission)
			assert.Equal(t, vo.Permission(tt.input), *permission)
			assert.Equal(t, tt.input, permission.String())
		})
	}
}

func TestPermission_FailureCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "empty string",
			input: "",
		},
		{
			name:  "no colon separator",
			input: "userslist",
		},
		{
			name:  "missing resource part",
			input: ":action",
		},
		{
			name:  "missing action part",
			input: "resource:",
		},
		{
			name:  "too long: 129 chars",
			input: strings.Repeat("a", 64) + ":" + strings.Repeat("b", 64), // 64 + 1 + 64 = 129
		},
		{
			name:  "multiple colons rejected",
			input: "users:list:read",
		},
		{
			name:  "colon only",
			input: ":",
		},
		{
			name:  "leading space treated as invalid (no TrimSpace)",
			input: " users:list",
		},
		{
			name:  "trailing space treated as invalid (no TrimSpace)",
			input: "users:list ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			permission, err := vo.NewPermission(tt.input)

			require.Error(t, err)
			assert.Nil(t, permission)

			var voErr vo.Error
			require.ErrorAs(t, err, &voErr)
			assert.Equal(t, vo.ValidationErrorCode, voErr.Code())
		})
	}
}
