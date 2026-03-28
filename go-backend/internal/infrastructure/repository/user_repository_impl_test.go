//go:build integration

package repository_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Haya372/web-app-template/go-backend/internal/domain/entity"
	domain_repository "github.com/Haya372/web-app-template/go-backend/internal/domain/entity/repository"
	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreate_HappyCase(t *testing.T) {
	target := repository.NewUserRepository(testDb.DbManager())
	tests := []struct {
		name string
		user entity.User
	}{
		{
			name: "Create Success",
			user: entity.ReconstructUser(
				uuid.New(),
				"test@example.com",
				[]byte("password"),
				"Test User",
				vo.UserStatusActive,
				time.Date(2026, 1, 18, 0, 0, 0, 0, time.UTC),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			user, err := target.Create(ctx, tt.user)

			assert.Nil(t, err)
			assert.NotNil(t, user)
			assert.Equal(t, user.ID(), tt.user.ID())
			assert.Equal(t, user.Email(), tt.user.Email())
			assert.Equal(t, user.PasswordHash(), tt.user.PasswordHash())
			assert.Equal(t, user.Name(), tt.user.Name())
			assert.Equal(t, user.CreatedAt(), tt.user.CreatedAt())
			assert.Equal(t, user.Status(), tt.user.Status())
		})
	}

	testDb.Cleanup()
}

func TestCreate_ErrorCase(t *testing.T) {
	target := repository.NewUserRepository(testDb.DbManager())
	tests := []struct {
		name string
		user entity.User
	}{
		{
			name: "Create Error",
			user: entity.ReconstructUser(
				uuid.New(),
				"test@example.com",
				[]byte("password"),
				strings.Repeat("a", 257),
				vo.UserStatusActive,
				time.Date(2026, 1, 18, 0, 0, 0, 0, time.UTC),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			user, err := target.Create(ctx, tt.user)

			assert.NotNil(t, err)
			assert.Nil(t, user)
		})
	}

	testDb.Cleanup()
}

func TestCreate_DuplicateEmail(t *testing.T) {
	target := repository.NewUserRepository(testDb.DbManager())
	ctx := context.Background()

	first := entity.ReconstructUser(
		uuid.New(),
		"dup@example.com",
		[]byte("password"),
		"First User",
		vo.UserStatusActive,
		time.Date(2026, 1, 18, 0, 0, 0, 0, time.UTC),
	)
	_, err := target.Create(ctx, first)
	if err != nil {
		t.Fatalf("failed to create first user: %v", err)
	}

	second := entity.ReconstructUser(
		uuid.New(),
		"dup@example.com",
		[]byte("other-password"),
		"Second User",
		vo.UserStatusActive,
		time.Date(2026, 1, 18, 0, 0, 0, 0, time.UTC),
	)
	user, err := target.Create(ctx, second)

	assert.Nil(t, user)
	assert.Error(t, err)

	var domainErr vo.Error
	assert.True(t, errors.As(err, &domainErr), "error should be a vo.Error")
	assert.Equal(t, vo.DuplicateEmailErrorCode, domainErr.Code())
	assert.Equal(t, 409, domainErr.Status())

	testDb.Cleanup()
}

func TestFindByEmail_HappyCase(t *testing.T) {
	seedUser := entity.ReconstructUser(
		uuid.New(),
		"test@example.com",
		[]byte("password"),
		"Test User",
		vo.UserStatusActive,
		time.Date(2026, 1, 18, 0, 0, 0, 0, time.UTC),
	)
	target := repository.NewUserRepository(testDb.DbManager())

	_, err := target.Create(context.Background(), seedUser)
	if err != nil {
		assert.Failf(t, "failed to create seed user", "err=%v", err)
	}

	tests := []struct {
		name   string
		email  string
		expect entity.User
	}{
		{
			name:   "Found user",
			email:  "test@example.com",
			expect: seedUser,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			user, err := target.FindByEmail(ctx, tt.email)

			assert.Nil(t, err)
			assert.Equal(t, user, tt.expect)
		})
	}

	testDb.Cleanup()
}

func TestFindByEmail_ErrorCase(t *testing.T) {
	target := repository.NewUserRepository(testDb.DbManager())

	tests := []struct {
		name  string
		email string
	}{
		{
			name:  "Notfound user",
			email: "notfound@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			user, err := target.FindByEmail(ctx, tt.email)

			assert.Error(t, err)
			assert.True(t, errors.Is(err, domain_repository.ErrUserNotFound))
			assert.Nil(t, user)
		})
	}

	testDb.Cleanup()
}
