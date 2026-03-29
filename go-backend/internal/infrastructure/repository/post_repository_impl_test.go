//go:build integration

package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/Haya372/web-app-template/go-backend/internal/domain/entity"
	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func seedUser(t *testing.T) entity.User {
	t.Helper()

	userRepo := repository.NewUserRepository(testDb.DbManager())
	user := entity.ReconstructUser(
		uuid.New(),
		"posttest@example.com",
		[]byte("password"),
		"Post Test User",
		vo.UserStatusActive,
		time.Date(2026, 3, 8, 0, 0, 0, 0, time.UTC),
	)

	created, err := userRepo.Create(context.Background(), user)
	require.NoError(t, err)

	return created
}

func TestCreatePost_HappyCase(t *testing.T) {
	user := seedUser(t)
	target := repository.NewPostRepository(testDb.DbManager())

	tests := []struct {
		name string
		post entity.Post
	}{
		{
			name: "Create post success",
			post: entity.ReconstructPost(
				uuid.New(),
				user.ID(),
				"Hello, world!",
				time.Date(2026, 3, 8, 12, 0, 0, 0, time.UTC),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			created, err := target.Create(ctx, tt.post)

			assert.NoError(t, err)
			assert.NotNil(t, created)
			assert.Equal(t, tt.post.ID(), created.ID())
			assert.Equal(t, tt.post.UserID(), created.UserID())
			assert.Equal(t, tt.post.Content(), created.Content())
		})
	}

	testDb.Cleanup()
}

func TestCreatePost_ForeignKeyViolation(t *testing.T) {
	defer func() { require.NoError(t, testDb.Cleanup()) }()

	// Use a random UUID that has no corresponding user row to trigger a FK violation.
	nonExistentUserID := uuid.New()
	post := entity.ReconstructPost(
		uuid.New(),
		nonExistentUserID,
		"this should fail",
		time.Date(2026, 3, 8, 12, 0, 0, 0, time.UTC),
	)

	target := repository.NewPostRepository(testDb.DbManager())
	_, err := target.Create(context.Background(), post)

	require.Error(t, err)
}

func TestCreatePost_DuplicateID(t *testing.T) {
	defer func() { require.NoError(t, testDb.Cleanup()) }()

	user := seedUser(t)
	fixedID := uuid.New()

	post := entity.ReconstructPost(
		fixedID,
		user.ID(),
		"first post",
		time.Date(2026, 3, 8, 12, 0, 0, 0, time.UTC),
	)

	target := repository.NewPostRepository(testDb.DbManager())

	_, err := target.Create(context.Background(), post)
	require.NoError(t, err)

	duplicate := entity.ReconstructPost(
		fixedID,
		user.ID(),
		"duplicate post",
		time.Date(2026, 3, 8, 13, 0, 0, 0, time.UTC),
	)
	_, err = target.Create(context.Background(), duplicate)

	require.Error(t, err)
}
