//go:build integration

package repository_test

import (
	"context"
	"testing"

	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Seeded role IDs from db/seeds/master/seed.sql.
const (
	adminRoleID  = "00000000-0000-0000-0000-000000000001"
	viewerRoleID = "00000000-0000-0000-0000-000000000002"
)

func assignRole(t *testing.T, userId, roleId string) {
	t.Helper()

	_, err := testDb.Pool().Exec(
		context.Background(),
		"INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		userId, roleId,
	)
	require.NoError(t, err)
}

func TestUserPermissionRepository_FindByUserId_WithAdminRole(t *testing.T) {
	defer func() { require.NoError(t, testDb.Cleanup()) }()

	u := seedUser(t)
	assignRole(t, u.Id().String(), adminRoleID)

	r := repository.NewUserPermissionRepository(testDb.DbManager())
	agg, err := r.FindByUserId(context.Background(), u.Id())

	require.NoError(t, err)
	require.NotNil(t, agg)

	assert.Equal(t, u.Id(), agg.UserId)
	assert.Equal(t, u.Email(), agg.User.Email())
	assert.Equal(t, u.Name(), agg.User.Name())
	assert.True(t, agg.HasPermission(vo.PermissionUsersList))
}

func TestUserPermissionRepository_FindByUserId_WithViewerRole(t *testing.T) {
	defer func() { require.NoError(t, testDb.Cleanup()) }()

	u := seedUser(t)
	assignRole(t, u.Id().String(), viewerRoleID)

	r := repository.NewUserPermissionRepository(testDb.DbManager())
	agg, err := r.FindByUserId(context.Background(), u.Id())

	require.NoError(t, err)
	assert.True(t, agg.HasPermission(vo.PermissionUsersList))
}

func TestUserPermissionRepository_FindByUserId_NoRole(t *testing.T) {
	defer func() { require.NoError(t, testDb.Cleanup()) }()

	u := seedUser(t)

	r := repository.NewUserPermissionRepository(testDb.DbManager())
	agg, err := r.FindByUserId(context.Background(), u.Id())

	require.NoError(t, err)
	require.NotNil(t, agg)
	assert.False(t, agg.HasPermission(vo.PermissionUsersList))
	assert.Empty(t, agg.Permissions)
}

func TestUserPermissionRepository_FindByUserId_MultipleRoles(t *testing.T) {
	defer func() { require.NoError(t, testDb.Cleanup()) }()

	u := seedUser(t)
	assignRole(t, u.Id().String(), adminRoleID)
	assignRole(t, u.Id().String(), viewerRoleID)

	r := repository.NewUserPermissionRepository(testDb.DbManager())
	agg, err := r.FindByUserId(context.Background(), u.Id())

	require.NoError(t, err)
	assert.True(t, agg.HasPermission(vo.PermissionUsersList))
}

func TestUserPermissionRepository_FindByUserId_UserNotFound(t *testing.T) {
	defer func() { require.NoError(t, testDb.Cleanup()) }()

	r := repository.NewUserPermissionRepository(testDb.DbManager())
	agg, err := r.FindByUserId(context.Background(), uuid.New())

	assert.Error(t, err)
	assert.Nil(t, agg)
}
