package vo

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

// Permission represents a fine-grained access right using the "<resource>:<action>" format.
type Permission string

const (
	PermissionUsersList   Permission = "users:list"
	PermissionUsersCreate Permission = "users:create"

	// maxPermissionLength corresponds to the DB schema: permissions.code varchar(128).
	maxPermissionLength = 128
)

var errIllegalPermission = errors.New("illegal permission")

// NewPermission validates raw and returns a Permission value object.
// Permission is an internal code value; no whitespace trimming is applied
// (leading/trailing spaces are treated as invalid).
// The format must be "<resource>:<action>" where both parts are non-empty.
func NewPermission(raw string) (*Permission, error) {
	if raw == "" {
		return nil, NewValidationError("permission is required", nil, errIllegalPermission)
	}

	// Permission is an internal code value; whitespace is not trimmed and is treated as invalid.
	if strings.TrimSpace(raw) != raw {
		return nil, NewValidationError(
			"permission must be in format 'resource:action'",
			nil,
			errIllegalPermission,
		)
	}

	if utf8.RuneCountInString(raw) > maxPermissionLength {
		return nil, NewValidationError(
			fmt.Sprintf("permission must be at most %d characters long", maxPermissionLength),
			map[string]any{"max_length": maxPermissionLength},
			errIllegalPermission,
		)
	}

	// "<resource>:<action>" format: exactly one colon, non-empty parts on both sides.
	if strings.Count(raw, ":") != 1 {
		return nil, NewValidationError(
			"permission must be in format 'resource:action'",
			nil,
			errIllegalPermission,
		)
	}

	const colonParts = 2

	parts := strings.SplitN(raw, ":", colonParts)
	if parts[0] == "" || parts[1] == "" {
		return nil, NewValidationError(
			"permission must be in format 'resource:action'",
			nil,
			errIllegalPermission,
		)
	}

	permission := Permission(raw)

	return &permission, nil
}

func (p Permission) String() string {
	return string(p)
}
