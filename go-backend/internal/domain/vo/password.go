package vo

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

type Password string

const (
	// minimumPasswordLength is measured in Unicode code points (runes) for user-visible character count.
	minimumPasswordLength = 8
	// maximumPasswordLength is measured in bytes because bcrypt truncates input at 72 bytes;
	// capping here prevents silent truncation that would make 73+ byte passwords indistinguishable.
	maximumPasswordLength = 72
)

var (
	errIllegalPasswordLength = errors.New("illegal password length")
)

func NewPassword(raw string) (*Password, error) {
	// minimum is measured in runes (user-visible characters)
	if utf8.RuneCountInString(raw) < minimumPasswordLength {
		return nil, NewValidationError(
			fmt.Sprintf("password must be at least %d characters long", minimumPasswordLength), nil, errIllegalPasswordLength)
	}

	// maximum is measured in bytes because bcrypt truncates at 72 bytes
	if len(raw) > maximumPasswordLength {
		return nil, NewValidationError(
			fmt.Sprintf("password must be at most %d bytes long", maximumPasswordLength),
			map[string]any{"max_length": maximumPasswordLength},
			errIllegalPasswordLength)
	}

	password := Password(raw)

	return &password, nil
}

// String implements fmt.Stringer and returns a masked value to prevent
// accidental password exposure in logs or error messages.
func (p Password) String() string {
	return "[REDACTED]"
}
