//go:build integration

package query_test

import (
	"context"
	"testing"
	"time"

	"github.com/Haya372/web-app-template/go-backend/internal/domain/entity"
	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/query"
	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func seedPostUser(t *testing.T, email string) entity.User {
	t.Helper()

	u := entity.ReconstructUser(
		uuid.New(),
		email,
		[]byte("hash"),
		"Post User",
		vo.UserStatusActive,
		time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	)
	repo := repository.NewUserRepository(testDb.DbManager())
	created, err := repo.Create(context.Background(), u)
	require.NoError(t, err)
	return created
}

func seedPost(t *testing.T, userID uuid.UUID, content string, createdAt time.Time) entity.Post {
	t.Helper()

	p := entity.ReconstructPost(uuid.New(), userID, content, createdAt)
	repo := repository.NewPostRepository(testDb.DbManager())
	created, err := repo.Create(context.Background(), p)
	require.NoError(t, err)
	return created
}

func TestPostQueryService_FindAll_ReturnsPostsOrderedByCreatedAtDesc(t *testing.T) {
	defer func() { require.NoError(t, testDb.Cleanup()) }()

	user := seedPostUser(t, "findall@example.com")
	older := seedPost(t, user.ID(), "older post", time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC))
	newer := seedPost(t, user.ID(), "newer post", time.Date(2026, 1, 2, 10, 0, 0, 0, time.UTC))

	svc := query.NewPostQueryService(testDb.DbManager())
	posts, total, err := svc.FindAll(context.Background(), 10, 0)

	require.NoError(t, err)
	assert.Equal(t, 2, total)
	require.Len(t, posts, 2)
	// ORDER BY created_at DESC: newer first
	assert.Equal(t, newer.ID(), posts[0].ID)
	assert.Equal(t, older.ID(), posts[1].ID)
}

func TestPostQueryService_FindAll_EmptyTable(t *testing.T) {
	defer func() { require.NoError(t, testDb.Cleanup()) }()

	svc := query.NewPostQueryService(testDb.DbManager())
	posts, total, err := svc.FindAll(context.Background(), 10, 0)

	require.NoError(t, err)
	assert.Equal(t, 0, total)
	assert.NotNil(t, posts)
	assert.Empty(t, posts)
}

func TestPostQueryService_FindAll_Pagination(t *testing.T) {
	defer func() { require.NoError(t, testDb.Cleanup()) }()

	user := seedPostUser(t, "pagination@example.com")
	for i := range 5 {
		seedPost(t, user.ID(), "post content", time.Date(2026, 1, i+1, 0, 0, 0, 0, time.UTC))
	}

	svc := query.NewPostQueryService(testDb.DbManager())

	page1, total, err := svc.FindAll(context.Background(), 2, 0)
	require.NoError(t, err)
	assert.Equal(t, 5, total)
	assert.Len(t, page1, 2)

	page2, _, err := svc.FindAll(context.Background(), 2, 2)
	require.NoError(t, err)
	assert.Len(t, page2, 2)

	page3, _, err := svc.FindAll(context.Background(), 2, 4)
	require.NoError(t, err)
	assert.Len(t, page3, 1)
}

func TestPostQueryService_FindAll_BeyondOffsetReturnsEmpty(t *testing.T) {
	defer func() { require.NoError(t, testDb.Cleanup()) }()

	user := seedPostUser(t, "beyond@example.com")
	seedPost(t, user.ID(), "only post", time.Now().UTC())

	svc := query.NewPostQueryService(testDb.DbManager())
	posts, total, err := svc.FindAll(context.Background(), 10, 100)

	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.NotNil(t, posts)
	assert.Empty(t, posts)
}

func TestPostQueryService_FindAll_MapsFieldsCorrectly(t *testing.T) {
	defer func() { require.NoError(t, testDb.Cleanup()) }()

	user := seedPostUser(t, "mapping@example.com")
	created := seedPost(t, user.ID(), "hello mapping", time.Date(2026, 3, 1, 12, 0, 0, 0, time.UTC))

	svc := query.NewPostQueryService(testDb.DbManager())
	posts, total, err := svc.FindAll(context.Background(), 10, 0)

	require.NoError(t, err)
	assert.Equal(t, 1, total)
	require.Len(t, posts, 1)
	assert.Equal(t, created.ID(), posts[0].ID)
	assert.Equal(t, user.ID(), posts[0].UserID)
	assert.Equal(t, "hello mapping", posts[0].Content)
	assert.Equal(t, time.Date(2026, 3, 1, 12, 0, 0, 0, time.UTC), posts[0].CreatedAt)
}

func TestPostQueryService_FindAll_TotalAndCountAreConsistent(t *testing.T) {
	defer func() { require.NoError(t, testDb.Cleanup()) }()

	user := seedPostUser(t, "consistent@example.com")
	for i := range 7 {
		seedPost(t, user.ID(), "post", time.Date(2026, 1, i+1, 0, 0, 0, 0, time.UTC))
	}

	svc := query.NewPostQueryService(testDb.DbManager())

	// Fetch only 3 rows but total should be 7
	posts, total, err := svc.FindAll(context.Background(), 3, 0)
	require.NoError(t, err)
	assert.Equal(t, 7, total)
	assert.Len(t, posts, 3)
}
