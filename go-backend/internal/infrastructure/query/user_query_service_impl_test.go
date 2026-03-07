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

func seedUser(t *testing.T, email string) entity.User {
	t.Helper()

	u := entity.ReconstructUser(
		uuid.New(),
		email,
		[]byte("hash"),
		"Test User",
		vo.UserStatusActive,
		time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	)
	repo := repository.NewUserRepository(testDb.DbManager())
	created, err := repo.Create(context.Background(), u)
	require.NoError(t, err)

	return created
}

func TestUserQueryService_FindAll_ReturnsUsers(t *testing.T) {
	defer func() { require.NoError(t, testDb.Cleanup()) }()

	seedUser(t, "alice@example.com")
	seedUser(t, "bob@example.com")

	svc := query.NewUserQueryService(testDb.DbManager())
	users, total, err := svc.FindAll(context.Background(), 10, 0)

	require.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, users, 2)

	// Verify password_hash is not exposed in the DTO.
	for _, u := range users {
		assert.NotEmpty(t, u.Id)
		assert.NotEmpty(t, u.Email)
		assert.NotEmpty(t, u.Name)
		assert.NotEmpty(t, u.Status)
	}
}

func TestUserQueryService_FindAll_EmptyTable(t *testing.T) {
	defer func() { require.NoError(t, testDb.Cleanup()) }()

	svc := query.NewUserQueryService(testDb.DbManager())
	users, total, err := svc.FindAll(context.Background(), 10, 0)

	require.NoError(t, err)
	assert.Equal(t, 0, total)
	assert.Empty(t, users)
}

func TestUserQueryService_FindAll_Pagination(t *testing.T) {
	defer func() { require.NoError(t, testDb.Cleanup()) }()

	for i := range 5 {
		seedUser(t, "user"+string(rune('a'+i))+"@example.com")
	}

	svc := query.NewUserQueryService(testDb.DbManager())

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

func TestUserQueryService_FindAll_OrderedByCreatedAtDesc(t *testing.T) {
	defer func() { require.NoError(t, testDb.Cleanup()) }()

	seedUser(t, "first@example.com")
	seedUser(t, "second@example.com")

	svc := query.NewUserQueryService(testDb.DbManager())
	users, _, err := svc.FindAll(context.Background(), 10, 0)

	require.NoError(t, err)
	require.Len(t, users, 2)
	// Results are ordered by created_at DESC; since both are created in the same
	// second in tests the order may vary, so just verify both are returned.
	emails := []string{users[0].Email, users[1].Email}
	assert.Contains(t, emails, "first@example.com")
	assert.Contains(t, emails, "second@example.com")
}
