package vo

import "errors"

type Password string

var (
	illegalPasswordLength = errors.New("illegal password length")
)

func NewPassword(raw string) (*Password, error) {
	if len(raw) < 8 {
		return nil, NewValidationError("password must be at least 8 characters long", nil, illegalPasswordLength)
	}

	password := Password(raw)

	return &password, nil
}
