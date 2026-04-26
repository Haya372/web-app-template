package vo

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"unicode/utf8"
)

type Email string

const (
	maxEmailLength = 256
)

var (
	errIllegalEmail = errors.New("illegal email")
)

func NewEmail(raw string) (*Email, error) {
	trimmed := strings.TrimSpace(raw)

	if trimmed == "" {
		return nil, NewValidationError("email is required", nil, errIllegalEmail)
	}

	if utf8.RuneCountInString(trimmed) > maxEmailLength {
		return nil, NewValidationError(
			fmt.Sprintf("email must be at most %d characters long", maxEmailLength),
			map[string]any{"max_length": maxEmailLength},
			errIllegalEmail,
		)
	}

	// addr.Address is the normalised bare address extracted by the parser.
	// Requiring equality with the trimmed input rejects display-name forms,
	// quoted local parts, and comments — only plain addr@domain strings pass.
	addr, err := mail.ParseAddress(trimmed)
	if err != nil || addr.Address != trimmed {
		return nil, NewValidationError("invalid email format", nil, errIllegalEmail)
	}

	// Domain is case-insensitive per RFC 5321 §2.4; normalise to lowercase.
	// The local part is case-preserving.
	// Use addr.Address (the parsed, validated form) as the source to make the
	// data-flow dependency explicit and resilient to future changes.
	const atSignParts = 2

	parts := strings.SplitN(addr.Address, "@", atSignParts)
	normalised := parts[0] + "@" + strings.ToLower(parts[1])
	email := Email(normalised)

	return &email, nil
}

func (e Email) String() string {
	return string(e)
}
