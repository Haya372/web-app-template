package vo_test

import (
	"testing"

	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContent_HappyCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "non-empty string creates valid Content VO",
			input: "Hello, world!",
		},
		{
			name:  "whitespace-only string is treated as valid",
			input: "   ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := vo.NewContent(tt.input)

			require.NoError(t, err)
			assert.Equal(t, vo.Content(tt.input), *content)
		})
	}
}

func TestContent_FailureCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "empty string returns ValidationError",
			input: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := vo.NewContent(tt.input)

			require.Error(t, err)
			assert.Nil(t, content)

			// Verify the error is a ValidationError (has the expected code)
			var voErr vo.Error
			require.ErrorAs(t, err, &voErr)
			assert.Equal(t, vo.ValidationErrorCode, voErr.Code())
		})
	}
}
