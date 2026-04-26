package vo_test

import (
	"strings"
	"testing"

	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmail_HappyCase(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantEmail string
	}{
		{
			name:      "valid email",
			input:     "user@example.com",
			wantEmail: "user@example.com",
		},
		{
			name:      "email with plus tag",
			input:     "user+tag@sub.domain.com",
			wantEmail: "user+tag@sub.domain.com",
		},
		{
			name:      "boundary: exactly 256-char valid email",
			input:     strings.Repeat("a", 244) + "@example.com", // 244 + 12 = 256 chars
			wantEmail: strings.Repeat("a", 244) + "@example.com",
		},
		{
			name:      "whitespace-padded email is trimmed",
			input:     " user@example.com ",
			wantEmail: "user@example.com",
		},
		{
			name:      "subdomain with hyphens and numbers",
			input:     "user@mail-server.co.jp",
			wantEmail: "user@mail-server.co.jp",
		},
		{
			name:      "uppercase domain is normalised to lowercase",
			input:     "user@EXAMPLE.COM",
			wantEmail: "user@example.com",
		},
		{
			name:      "mixed-case local part is preserved",
			input:     "User@example.com",
			wantEmail: "User@example.com",
		},
		{
			name:      "whitespace trimming and domain normalisation compose correctly",
			input:     " user@EXAMPLE.COM ",
			wantEmail: "user@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := vo.NewEmail(tt.input)

			require.NoError(t, err)
			assert.NotNil(t, email)
			assert.Equal(t, vo.Email(tt.wantEmail), *email)
			assert.Equal(t, tt.wantEmail, email.String())
		})
	}
}

func TestEmail_FailureCase(t *testing.T) {
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
			name:  "no at sign",
			input: "notanemail",
		},
		{
			name:  "missing domain",
			input: "user@",
		},
		{
			name:  "missing local part",
			input: "@domain.com",
		},
		{
			name:  "display name format rejected",
			input: `"Name" <user@example.com>`,
		},
		{
			name:  "too long: 257 chars",
			input: strings.Repeat("a", 245) + "@example.com", // 245 + 12 = 257 chars
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := vo.NewEmail(tt.input)

			require.Error(t, err)
			assert.Nil(t, email)

			var voErr vo.Error
			require.ErrorAs(t, err, &voErr)
			assert.Equal(t, vo.ValidationErrorCode, voErr.Code())
		})
	}
}
