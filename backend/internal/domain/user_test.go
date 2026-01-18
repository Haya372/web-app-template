package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
			email:     "test@example.com",
			password:  "password",
			name:      "Test",
			createdAt: time.Date(2026, 1, 18, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			user, err := NewUser(tt.email, tt.password, tt.name, tt.createdAt)

			assert.Nil(t, err)
			assert.Equal(t, user.Name(), tt.name)
			assert.Equal(t, user.Email(), tt.email)
			assert.Equal(t, user.CreatedAt(), tt.createdAt)

			err = bcrypt.CompareHashAndPassword(user.PasswordHash(), []byte(tt.password))
			assert.Nil(t, err)
		})
	}
}
