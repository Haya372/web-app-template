package vo_test

import (
	"strings"
	"testing"

	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestName_HappyCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantName string
	}{
		{
			name:     "valid ASCII name",
			input:    "Alice",
			wantName: "Alice",
		},
		{
			name:     "Japanese name",
			input:    "\u5c71\u7530 \u592a\u90ce",
			wantName: "\u5c71\u7530 \u592a\u90ce",
		},
		{
			name:     "name with leading/trailing whitespace is trimmed",
			input:    " Bob ",
			wantName: "Bob",
		},
		{
			name:     "boundary: exactly 256-char name",
			input:    strings.Repeat("a", 256),
			wantName: strings.Repeat("a", 256),
		},
		{
			name:     "boundary: exactly 256-char Unicode name",
			input:    strings.Repeat("あ", 256),
			wantName: strings.Repeat("あ", 256),
		},
		{
			name:     "trim brings 258-char input to 256 chars and succeeds",
			input:    " " + strings.Repeat("a", 256) + " ",
			wantName: strings.Repeat("a", 256),
		},
		{
			name:     "4-byte emoji runes counted correctly at boundary",
			input:    strings.Repeat("😀", 256),
			wantName: strings.Repeat("😀", 256),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, err := vo.NewName(tt.input)

			require.NoError(t, err)
			assert.NotNil(t, name)
			assert.Equal(t, vo.Name(tt.wantName), *name)
			assert.Equal(t, tt.wantName, name.String())
		})
	}
}

func TestName_FailureCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "empty string",
			input: "",
		},
		{
			name:  "whitespace only",
			input: "   ",
		},
		{
			name:  "too long ASCII: 257 chars",
			input: strings.Repeat("a", 257),
		},
		{
			name:  "too long Unicode: 257 runes",
			input: strings.Repeat("あ", 257),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, err := vo.NewName(tt.input)

			require.Error(t, err)
			assert.Nil(t, name)

			var voErr vo.Error
			require.ErrorAs(t, err, &voErr)
			assert.Equal(t, vo.ValidationErrorCode, voErr.Code())
		})
	}
}
