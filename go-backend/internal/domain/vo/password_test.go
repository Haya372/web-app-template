package vo_test

import (
	"strings"
	"testing"

	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPassword_HappyCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "boundary: exactly 8 chars (minimum)",
			input: "pass1234",
		},
		{
			name:  "boundary: exactly 72 bytes (maximum ASCII)",
			input: strings.Repeat("a", 72),
		},
		{
			name:  "boundary: 8 multi-byte runes satisfies minimum",
			input: strings.Repeat("あ", 8), // 8 runes = 24 bytes, passes rune-based min check
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password, err := vo.NewPassword(tt.input)

			require.NoError(t, err)
			assert.Equal(t, *password, vo.Password(tt.input))
		})
	}
}

func TestPassword_FailureCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "password length under 8 characters",
			input: "passwor",
		},
		{
			name:  "too long: 73 bytes exceeds maximum",
			input: strings.Repeat("a", 73),
		},
		{
			name:  "too long: multi-byte input exceeds 72 bytes even if fewer runes",
			input: strings.Repeat("あ", 25), // 25 runes = 75 bytes, exceeds byte limit
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password, err := vo.NewPassword(tt.input)

			require.Error(t, err)
			assert.Nil(t, password)

			var voErr vo.Error
			require.ErrorAs(t, err, &voErr)
			assert.Equal(t, vo.ValidationErrorCode, voErr.Code())
		})
	}
}

func TestPassword_MaxLengthErrorDetails(t *testing.T) {
	_, err := vo.NewPassword(strings.Repeat("a", 73))

	require.Error(t, err)

	var voErr vo.Error
	require.ErrorAs(t, err, &voErr)
	assert.Equal(t, map[string]any{"max_length": 72}, voErr.Details())
}

func TestPassword_String(t *testing.T) {
	password, err := vo.NewPassword("password")
	require.NoError(t, err)
	assert.Equal(t, "[REDACTED]", password.String())
}
