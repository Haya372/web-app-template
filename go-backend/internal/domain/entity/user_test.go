package entity_test

import (
	"testing"
	"time"

	"github.com/Haya372/web-app-template/go-backend/internal/domain/entity"
	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestUser_HappyCase(t *testing.T) {
	tests := []struct {
		testName  string
		email     string
		password  string
		name      string
		createdAt time.Time
	}{
		{
			testName:  "Success to create User",
			email:     "test@example.com",
			password:  "password",
			name:      "Test",
			createdAt: time.Date(2026, 1, 18, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			user, err := entity.NewUser(tt.email, tt.password, tt.name, tt.createdAt)

			require.NoError(t, err)
			assert.Equal(t, user.Name(), tt.name)
			assert.Equal(t, user.Email(), tt.email)
			assert.Equal(t, user.CreatedAt(), tt.createdAt)
			assert.Equal(t, vo.UserStatusActive, user.Status())

			err = bcrypt.CompareHashAndPassword(user.PasswordHash(), []byte(tt.password))
			require.NoError(t, err)
		})
	}
}

func TestUser_FailureCase(t *testing.T) {
	tests := []struct {
		testName  string
		email     string
		password  string
		name      string
		createdAt time.Time
	}{
		{
			testName:  "password length under 8 characters",
			email:     "test@example.com",
			password:  "passwor",
			name:      "Test",
			createdAt: time.Date(2026, 1, 18, 0, 0, 0, 0, time.UTC),
		},
		{
			testName:  "empty name",
			email:     "test@example.com",
			password:  "password",
			name:      "",
			createdAt: time.Date(2026, 1, 18, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			user, err := entity.NewUser(tt.email, tt.password, tt.name, tt.createdAt)

			require.Error(t, err)
			assert.Nil(t, user)
		})
	}
}

func TestUser_UpdateStatus(t *testing.T) {
	origin := time.Date(2026, 1, 18, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name        string
		current     vo.UserStatus
		target      vo.UserStatus
		shouldError bool
	}{
		{
			name:    "active to frozen",
			current: vo.UserStatusActive,
			target:  vo.UserStatusFrozen,
		},
		{
			name:    "frozen to active",
			current: vo.UserStatusFrozen,
			target:  vo.UserStatusActive,
		},
		{
			name:    "active to deleted",
			current: vo.UserStatusActive,
			target:  vo.UserStatusDeleted,
		},
		{
			name:    "frozen to deleted",
			current: vo.UserStatusFrozen,
			target:  vo.UserStatusDeleted,
		},
		{
			name:        "deleted to active should error",
			current:     vo.UserStatusDeleted,
			target:      vo.UserStatusActive,
			shouldError: true,
		},
		{
			name:        "deleted to frozen should error",
			current:     vo.UserStatusDeleted,
			target:      vo.UserStatusFrozen,
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := entity.ReconstructUser(
				uuid.New(),
				"test@example.com",
				[]byte("password"),
				"Test",
				tt.current,
				origin,
			)

			updated, err := user.UpdateStatus(tt.target)

			if tt.shouldError {
				require.Error(t, err)
				assert.Nil(t, updated)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.target, updated.Status())
			assert.Equal(t, user.Id(), updated.Id())
			assert.Equal(t, user.CreatedAt(), updated.CreatedAt())
		})
	}
}
