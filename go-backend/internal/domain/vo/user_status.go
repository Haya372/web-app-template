package vo

import (
	"errors"
	"strings"
)

type UserStatus string

const (
	UserStatusActive  UserStatus = "ACTIVE"
	UserStatusFrozen  UserStatus = "FROZEN"
	UserStatusDeleted UserStatus = "DELETED"
)

var (
	errInvalidUserStatus = errors.New("invalid user status")
)

func (s UserStatus) String() string {
	return string(s)
}

func (s UserStatus) IsActive() bool {
	return s == UserStatusActive
}

func (s UserStatus) IsFrozen() bool {
	return s == UserStatusFrozen
}

func (s UserStatus) IsDeleted() bool {
	return s == UserStatusDeleted
}

func UserStatusFromString(raw string) (UserStatus, error) {
	switch strings.ToUpper(raw) {
	case string(UserStatusActive):
		return UserStatusActive, nil
	case string(UserStatusFrozen):
		return UserStatusFrozen, nil
	case string(UserStatusDeleted):
		return UserStatusDeleted, nil
	default:
		return "", NewValidationError("invalid user status", map[string]any{
			"status": raw,
		}, errInvalidUserStatus)
	}
}
