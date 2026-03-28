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
