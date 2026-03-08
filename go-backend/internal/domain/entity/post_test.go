package entity_test

import (
	"testing"
	"time"

	"github.com/Haya372/web-app-template/go-backend/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPost_HappyCase(t *testing.T) {
	tests := []struct {
		name      string
		userId    uuid.UUID
		content   string
		createdAt time.Time
	}{
		{
			name:      "valid userId and non-empty content creates a Post",
			userId:    uuid.New(),
			content:   "This is a post content.",
			createdAt: time.Date(2026, 1, 18, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post, err := entity.NewPost(tt.userId, tt.content, tt.createdAt)

			require.NoError(t, err)
			assert.NotEqual(t, uuid.Nil, post.Id())
			assert.Equal(t, tt.userId, post.UserId())
			assert.Equal(t, tt.content, post.Content())
			assert.Equal(t, tt.createdAt, post.CreatedAt())
		})
	}
}

func TestNewPost_FailureCase(t *testing.T) {
	tests := []struct {
		name      string
		userId    uuid.UUID
		content   string
		createdAt time.Time
	}{
		{
			name:      "empty content returns error",
			userId:    uuid.New(),
			content:   "",
			createdAt: time.Date(2026, 1, 18, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post, err := entity.NewPost(tt.userId, tt.content, tt.createdAt)

			require.Error(t, err)
			assert.Nil(t, post)
		})
	}
}

func TestNewPost_GeneratesUuidV7(t *testing.T) {
	t.Run("generated id has UUID version 7", func(t *testing.T) {
		userId := uuid.New()
		post, err := entity.NewPost(userId, "some content", time.Now())

		require.NoError(t, err)
		assert.Equal(t, uuid.Version(7), post.Id().Version())
	})

	t.Run("consecutively generated ids are time-ordered", func(t *testing.T) {
		userId := uuid.New()
		post1, err1 := entity.NewPost(userId, "first post", time.Now())
		post2, err2 := entity.NewPost(userId, "second post", time.Now())

		require.NoError(t, err1)
		require.NoError(t, err2)
		// UUID v7 bytes are lexicographically sortable by creation time
		assert.True(t, post1.Id().String() < post2.Id().String() || post1.Id().String() == post2.Id().String(),
			"first post ID should be <= second post ID in lexicographic order")
	})
}

func TestReconstructPost_HappyCase(t *testing.T) {
	tests := []struct {
		name      string
		id        uuid.UUID
		userId    uuid.UUID
		content   string
		createdAt time.Time
	}{
		{
			name:      "reconstructs a Post with exact field values",
			id:        uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			userId:    uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			content:   "Reconstructed post content.",
			createdAt: time.Date(2026, 3, 8, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post := entity.ReconstructPost(tt.id, tt.userId, tt.content, tt.createdAt)

			assert.Equal(t, tt.id, post.Id())
			assert.Equal(t, tt.userId, post.UserId())
			assert.Equal(t, tt.content, post.Content())
			assert.Equal(t, tt.createdAt, post.CreatedAt())
		})
	}
}
