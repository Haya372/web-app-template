package vo_test

import (
	"strings"
	"testing"

	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContent_HappyCase(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantContent string
	}{
		{
			name:        "non-empty string creates valid Content VO",
			input:       "Hello, world!",
			wantContent: "Hello, world!",
		},
		{
			name:        "content with leading/trailing whitespace is trimmed",
			input:       "  hello  ",
			wantContent: "hello",
		},
		{
			name:        "boundary: exactly 10000-char content",
			input:       strings.Repeat("a", 10000),
			wantContent: strings.Repeat("a", 10000),
		},
		{
			name:        "boundary: exactly 10000-char Unicode content",
			input:       strings.Repeat("あ", 10000),
			wantContent: strings.Repeat("あ", 10000),
		},
		{
			name:        "trim brings over-limit-before-trim input to exactly 10000 runes and succeeds",
			input:       "  " + strings.Repeat("a", 10000) + "  ",
			wantContent: strings.Repeat("a", 10000),
		},
		{
			name:        "boundary: exactly 10000 four-byte emoji runes",
			input:       strings.Repeat("😀", 10000),
			wantContent: strings.Repeat("😀", 10000),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := vo.NewContent(tt.input)

			require.NoError(t, err)
			assert.Equal(t, vo.Content(tt.wantContent), *content)
			assert.Equal(t, tt.wantContent, content.String())
		})
	}
}

func TestContent_FailureCase(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantDetails map[string]any
	}{
		{
			name:        "empty string returns ValidationError",
			input:       "",
			wantDetails: nil,
		},
		{
			name:        "whitespace-only string returns ValidationError",
			input:       "   ",
			wantDetails: nil,
		},
		{
			name:        "too long: 10001 chars exceeds maximum",
			input:       strings.Repeat("a", 10001),
			wantDetails: map[string]any{"max_length": 10000},
		},
		{
			name:        "too long: 10001 Unicode runes exceeds maximum",
			input:       strings.Repeat("あ", 10001),
			wantDetails: map[string]any{"max_length": 10000},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := vo.NewContent(tt.input)

			require.Error(t, err)
			assert.Nil(t, content)

			var voErr vo.Error
			require.ErrorAs(t, err, &voErr)
			assert.Equal(t, vo.ValidationErrorCode, voErr.Code())
			assert.Equal(t, tt.wantDetails, voErr.Details())
		})
	}
}
