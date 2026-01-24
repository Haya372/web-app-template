package entity_test

import (
	"testing"
	"time"

	"github.com/Haya372/web-app-template/backend/internal/domain/entity"
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
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			user, err := entity.NewUser(tt.email, tt.password, tt.name, tt.createdAt)

			assert.Error(t, err)
			assert.Nil(t, user)
		})
	}
}
