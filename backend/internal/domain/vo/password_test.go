package vo_test

import (
	"testing"

	"github.com/Haya372/web-app-template/backend/internal/domain/vo"
	"github.com/stretchr/testify/assert"
)

func TestPassword_HappyCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "correct password",
			input: "password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password, err := vo.NewPassword(tt.input)

			assert.NoError(t, err)
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password, err := vo.NewPassword(tt.input)

			assert.Error(t, err)
			assert.Nil(t, password)
		})
	}
}
