package vo

import (
	"errors"
	"fmt"
)

type Password string

const (
	minimumPasswordLength = 8
)

var (
	errIllegalPasswordLength = errors.New("illegal password length")
)

func NewPassword(raw string) (*Password, error) {
	if len(raw) < minimumPasswordLength {
		return nil, NewValidationError(
			fmt.Sprintf("password must be at least %d characters long", minimumPasswordLength), nil, errIllegalPasswordLength)
	}

	password := Password(raw)

	return &password, nil
}
