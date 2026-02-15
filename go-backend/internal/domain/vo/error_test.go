package vo_test

import (
	"errors"
	"testing"

	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	"github.com/stretchr/testify/assert"
)

func TestNewValidationError(t *testing.T) {
	err := vo.NewValidationError("invalid input", map[string]any{"field": "email"}, errors.New("base"))

	var baseErr vo.Error
	if assert.ErrorAs(t, err, &baseErr) {
		assert.Equal(t, 400, baseErr.Status())
		assert.Equal(t, vo.ValidationErrorCode, baseErr.Code())
		assert.Equal(t, "invalid input", baseErr.Message())
		assert.Equal(t, map[string]any{"field": "email"}, baseErr.Details())
	}
}

func TestNewUnauthorizedError(t *testing.T) {
	err := vo.NewUnauthorizedError("invalid credential", nil, errors.New("base"))

	var baseErr vo.Error
	if assert.ErrorAs(t, err, &baseErr) {
		assert.Equal(t, 401, baseErr.Status())
		assert.Equal(t, vo.InvalidCredentialErrorCode, baseErr.Code())
		assert.Equal(t, "invalid credential", baseErr.Message())
		assert.Nil(t, baseErr.Details())
	}
}

func TestErrorCode_Title(t *testing.T) {
	tests := []struct {
		name     string
		code     vo.ErrorCode
		expected string
	}{
		{
			name:     "validation",
			code:     vo.ValidationErrorCode,
			expected: "validation error",
		},
		{
			name:     "invalid credential",
			code:     vo.InvalidCredentialErrorCode,
			expected: "invalid credential",
		},
		{
			name:     "internal",
			code:     vo.InternalErrorCode,
			expected: "internal server error",
		},
		{
			name:     "unknown",
			code:     vo.ErrorCode("UNKNOWN"),
			expected: "application error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.code.Title())
		})
	}
}
