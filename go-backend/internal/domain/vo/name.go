package vo

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

type Name string

const (
	maxNameLength = 256
)

var (
	errIllegalName = errors.New("illegal name")
)

func NewName(raw string) (*Name, error) {
	trimmed := strings.TrimSpace(raw)

	if trimmed == "" {
		return nil, NewValidationError("name is required", nil, errIllegalName)
	}

	if utf8.RuneCountInString(trimmed) > maxNameLength {
		return nil, NewValidationError(
			fmt.Sprintf("name must be at most %d characters long", maxNameLength),
			map[string]any{"max_length": maxNameLength},
			errIllegalName,
		)
	}

	name := Name(trimmed)

	return &name, nil
}

func (n Name) String() string {
	return string(n)
}
