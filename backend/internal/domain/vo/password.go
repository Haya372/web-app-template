package vo

import (
	"fmt"
)

type Password string

const (
	minimumPasswordLength = 8
)

func NewPassword(raw string) (*Password, error) {
	if len(raw) < minimumPasswordLength {
		return nil, NewValidationError(
			fmt.Sprintf("password must be at least %d characters long", minimumPasswordLength), nil, nil)
	}

	password := Password(raw)

	return &password, nil
}
