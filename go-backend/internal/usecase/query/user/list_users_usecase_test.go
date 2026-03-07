package user_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Haya372/web-app-template/go-backend/internal/domain/snapshot"
	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	"github.com/Haya372/web-app-template/go-backend/internal/usecase/query/user"
	mock_reader "github.com/Haya372/web-app-template/go-backend/test/mock/domain/snapshot/reader"
	mock_query "github.com/Haya372/web-app-template/go-backend/test/mock/usecase/query"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func newTestUseCase(
	t *testing.T,
	queryService user.UserQueryService,
	permReader *mock_reader.MockUserPermissionReader,
) user.ListUsersUseCase {
	t.Helper()

	return user.NewListUsersUseCase(queryService, permReader)
}

func withPermission(userID uuid.UUID, perms ...vo.Permission) *snapshot.UserPermissionSnapshot {
	return &snapshot.UserPermissionSnapshot{
		UserId:      userID,
		Permissions: perms,
	}
}

func TestListUsersUseCase_HappyCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	userID := uuid.New()

	now := time.Now().UTC()
	expectedUsers := []user.UserDto{
		{Id: uuid.New(), Name: "Alice", Email: "alice@example.com", Status: "ACTIVE", CreatedAt: now},
		{Id: uuid.New(), Name: "Bob", Email: "bob@example.com", Status: "ACTIVE", CreatedAt: now},
	}

	permReader := mock_reader.NewMockUserPermissionReader(ctrl)
	permReader.EXPECT().FindByUserId(gomock.Any(), userID).
		Return(withPermission(userID, vo.PermissionUsersList), nil).Times(1)

	queryService := mock_query.NewMockUserQueryService(ctrl)
	queryService.EXPECT().FindAll(gomock.Any(), 20, 0).Return(expectedUsers, 2, nil).Times(1)

	uc := newTestUseCase(t, queryService, permReader)
	output, err := uc.Execute(context.Background(), user.ListUsersInput{UserId: userID, Limit: 20, Offset: 0})

	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Equal(t, 2, output.Total)
	assert.Len(t, output.Users, 2)
	assert.Equal(t, "Alice", output.Users[0].Name)
}

func TestListUsersUseCase_Pagination(t *testing.T) {
	ctrl := gomock.NewController(t)
	userID := uuid.New()

	now := time.Now().UTC()
	expectedUsers := []user.UserDto{
		{Id: uuid.New(), Name: "Carol", Email: "carol@example.com", Status: "ACTIVE", CreatedAt: now},
	}

	permReader := mock_reader.NewMockUserPermissionReader(ctrl)
	permReader.EXPECT().FindByUserId(gomock.Any(), userID).
		Return(withPermission(userID, vo.PermissionUsersList), nil).Times(1)

	queryService := mock_query.NewMockUserQueryService(ctrl)
	queryService.EXPECT().FindAll(gomock.Any(), 5, 10).Return(expectedUsers, 42, nil).Times(1)

	uc := newTestUseCase(t, queryService, permReader)
	output, err := uc.Execute(context.Background(), user.ListUsersInput{UserId: userID, Limit: 5, Offset: 10})

	require.NoError(t, err)
	assert.Equal(t, 42, output.Total)
	assert.Len(t, output.Users, 1)
}

func TestListUsersUseCase_EmptyResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	userID := uuid.New()

	permReader := mock_reader.NewMockUserPermissionReader(ctrl)
	permReader.EXPECT().FindByUserId(gomock.Any(), userID).
		Return(withPermission(userID, vo.PermissionUsersList), nil).Times(1)

	queryService := mock_query.NewMockUserQueryService(ctrl)
	queryService.EXPECT().FindAll(gomock.Any(), 20, 0).Return(nil, 0, nil).Times(1)

	uc := newTestUseCase(t, queryService, permReader)
	output, err := uc.Execute(context.Background(), user.ListUsersInput{UserId: userID, Limit: 20, Offset: 0})

	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Equal(t, 0, output.Total)
	assert.Empty(t, output.Users)
}

func TestListUsersUseCase_InvalidLimit(t *testing.T) {
	tests := []struct {
		name  string
		limit int
	}{
		{"limit zero", 0},
		{"limit negative", -1},
		{"limit over max", 101},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			userID := uuid.New()

			permReader := mock_reader.NewMockUserPermissionReader(ctrl)
			permReader.EXPECT().FindByUserId(gomock.Any(), userID).
				Return(withPermission(userID, vo.PermissionUsersList), nil).Times(1)

			queryService := mock_query.NewMockUserQueryService(ctrl)
			// FindAll must not be called on invalid input
			queryService.EXPECT().FindAll(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

			uc := newTestUseCase(t, queryService, permReader)
			output, err := uc.Execute(
				context.Background(), user.ListUsersInput{UserId: userID, Limit: tt.limit, Offset: 0},
			)

			require.Error(t, err)
			assert.Nil(t, output)

			var voErr vo.Error
			require.ErrorAs(t, err, &voErr)
			assert.Equal(t, vo.ValidationErrorCode, voErr.Code())
		})
	}
}

func TestListUsersUseCase_InvalidOffset(t *testing.T) {
	ctrl := gomock.NewController(t)
	userID := uuid.New()

	permReader := mock_reader.NewMockUserPermissionReader(ctrl)
	permReader.EXPECT().FindByUserId(gomock.Any(), userID).
		Return(withPermission(userID, vo.PermissionUsersList), nil).Times(1)

	queryService := mock_query.NewMockUserQueryService(ctrl)
	queryService.EXPECT().FindAll(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	uc := newTestUseCase(t, queryService, permReader)
	output, err := uc.Execute(context.Background(), user.ListUsersInput{UserId: userID, Limit: 20, Offset: -1})

	require.Error(t, err)
	assert.Nil(t, output)

	var voErr vo.Error
	require.ErrorAs(t, err, &voErr)
	assert.Equal(t, vo.ValidationErrorCode, voErr.Code())
}

func TestListUsersUseCase_QueryServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	userID := uuid.New()

	permReader := mock_reader.NewMockUserPermissionReader(ctrl)
	permReader.EXPECT().FindByUserId(gomock.Any(), userID).
		Return(withPermission(userID, vo.PermissionUsersList), nil).Times(1)

	queryService := mock_query.NewMockUserQueryService(ctrl)
	queryService.EXPECT().FindAll(gomock.Any(), 20, 0).Return(nil, 0, errors.New("db error")).Times(1)

	uc := newTestUseCase(t, queryService, permReader)
	output, err := uc.Execute(context.Background(), user.ListUsersInput{UserId: userID, Limit: 20, Offset: 0})

	require.Error(t, err)
	assert.Nil(t, output)
}

func TestListUsersUseCase_ForbiddenWhenNoPermission(t *testing.T) {
	ctrl := gomock.NewController(t)
	userID := uuid.New()

	permReader := mock_reader.NewMockUserPermissionReader(ctrl)
	permReader.EXPECT().FindByUserId(gomock.Any(), userID).
		Return(withPermission(userID), nil).Times(1)

	queryService := mock_query.NewMockUserQueryService(ctrl)
	queryService.EXPECT().FindAll(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	uc := newTestUseCase(t, queryService, permReader)
	output, err := uc.Execute(context.Background(), user.ListUsersInput{UserId: userID, Limit: 20, Offset: 0})

	require.Error(t, err)
	assert.Nil(t, output)

	var voErr vo.Error
	require.ErrorAs(t, err, &voErr)
	assert.Equal(t, vo.ForbiddenErrorCode, voErr.Code())
	assert.Equal(t, 403, voErr.Status())
}

func TestListUsersUseCase_ForbiddenWhenPermissionReaderError(t *testing.T) {
	ctrl := gomock.NewController(t)
	userID := uuid.New()

	permReader := mock_reader.NewMockUserPermissionReader(ctrl)
	permReader.EXPECT().FindByUserId(gomock.Any(), userID).
		Return(nil, errors.New("reader error")).Times(1)

	queryService := mock_query.NewMockUserQueryService(ctrl)
	queryService.EXPECT().FindAll(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	uc := newTestUseCase(t, queryService, permReader)
	output, err := uc.Execute(context.Background(), user.ListUsersInput{UserId: userID, Limit: 20, Offset: 0})

	require.Error(t, err)
	assert.Nil(t, output)
}

func TestListUsersUseCase_AllowedWhenOneOfMultipleRolesHasPermission(t *testing.T) {
	ctrl := gomock.NewController(t)
	userID := uuid.New()

	now := time.Now().UTC()
	expectedUsers := []user.UserDto{
		{Id: uuid.New(), Name: "Dave", Email: "dave@example.com", Status: "ACTIVE", CreatedAt: now},
	}

	permReader := mock_reader.NewMockUserPermissionReader(ctrl)
	permReader.EXPECT().FindByUserId(gomock.Any(), userID).
		Return(withPermission(userID, vo.PermissionUsersCreate, vo.PermissionUsersList), nil).Times(1)

	queryService := mock_query.NewMockUserQueryService(ctrl)
	queryService.EXPECT().FindAll(gomock.Any(), 20, 0).Return(expectedUsers, 1, nil).Times(1)

	uc := newTestUseCase(t, queryService, permReader)
	output, err := uc.Execute(context.Background(), user.ListUsersInput{UserId: userID, Limit: 20, Offset: 0})

	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Equal(t, 1, output.Total)
}
