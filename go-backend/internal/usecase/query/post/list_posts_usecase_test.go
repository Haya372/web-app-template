package post_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	"github.com/Haya372/web-app-template/go-backend/internal/usecase/query/post"
	mock_query "github.com/Haya372/web-app-template/go-backend/test/mock/usecase/query"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func newTestUseCase(t *testing.T, queryService post.PostQueryService) post.ListPostsUseCase {
	t.Helper()

	return post.NewListPostsUseCase(queryService)
}

func TestListPostsUseCase_HappyCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	now := time.Now().UTC()

	expectedPosts := []post.PostDto{
		{ID: uuid.New(), UserID: uuid.New(), Content: "Hello World", CreatedAt: now},
		{ID: uuid.New(), UserID: uuid.New(), Content: "Second post", CreatedAt: now},
	}

	queryService := mock_query.NewMockPostQueryService(ctrl)
	queryService.EXPECT().FindAll(gomock.Any(), 20, 0).Return(expectedPosts, 2, nil).Times(1)

	uc := newTestUseCase(t, queryService)
	output, err := uc.Execute(context.Background(), post.ListPostsInput{Limit: 20, Offset: 0})

	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Equal(t, 2, output.Total)
	assert.Len(t, output.Posts, 2)
	assert.Equal(t, "Hello World", output.Posts[0].Content)
}

func TestListPostsUseCase_Pagination(t *testing.T) {
	ctrl := gomock.NewController(t)
	now := time.Now().UTC()

	expectedPosts := []post.PostDto{
		{ID: uuid.New(), UserID: uuid.New(), Content: "Third", CreatedAt: now},
	}

	queryService := mock_query.NewMockPostQueryService(ctrl)
	queryService.EXPECT().FindAll(gomock.Any(), 5, 10).Return(expectedPosts, 42, nil).Times(1)

	uc := newTestUseCase(t, queryService)
	output, err := uc.Execute(context.Background(), post.ListPostsInput{Limit: 5, Offset: 10})

	require.NoError(t, err)
	assert.Equal(t, 42, output.Total)
	assert.Len(t, output.Posts, 1)
}

func TestListPostsUseCase_EmptyResult(t *testing.T) {
	ctrl := gomock.NewController(t)

	queryService := mock_query.NewMockPostQueryService(ctrl)
	queryService.EXPECT().FindAll(gomock.Any(), 20, 0).Return(nil, 0, nil).Times(1)

	uc := newTestUseCase(t, queryService)
	output, err := uc.Execute(context.Background(), post.ListPostsInput{Limit: 20, Offset: 0})

	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Equal(t, 0, output.Total)
	assert.Empty(t, output.Posts)
}

func TestListPostsUseCase_BoundaryLimit(t *testing.T) {
	tests := []struct {
		name  string
		limit int
		valid bool
	}{
		{"limit min valid", 1, true},
		{"limit max valid", 100, true},
		{"limit zero", 0, false},
		{"limit negative", -1, false},
		{"limit over max", 101, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			queryService := mock_query.NewMockPostQueryService(ctrl)

			if tt.valid {
				queryService.EXPECT().FindAll(gomock.Any(), tt.limit, 0).Return([]post.PostDto{}, 0, nil).Times(1)
			} else {
				queryService.EXPECT().FindAll(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			}

			uc := newTestUseCase(t, queryService)
			output, err := uc.Execute(context.Background(), post.ListPostsInput{Limit: tt.limit, Offset: 0})

			if tt.valid {
				require.NoError(t, err)
				require.NotNil(t, output)
			} else {
				require.Error(t, err)
				assert.Nil(t, output)

				var voErr vo.Error
				require.ErrorAs(t, err, &voErr)
				assert.Equal(t, vo.ValidationErrorCode, voErr.Code())
			}
		})
	}
}

func TestListPostsUseCase_OffsetZeroIsValid(t *testing.T) {
	ctrl := gomock.NewController(t)

	queryService := mock_query.NewMockPostQueryService(ctrl)
	queryService.EXPECT().FindAll(gomock.Any(), 20, 0).Return([]post.PostDto{}, 0, nil).Times(1)

	uc := newTestUseCase(t, queryService)
	output, err := uc.Execute(context.Background(), post.ListPostsInput{Limit: 20, Offset: 0})

	require.NoError(t, err)
	require.NotNil(t, output)
}

func TestListPostsUseCase_InvalidOffset(t *testing.T) {
	ctrl := gomock.NewController(t)

	queryService := mock_query.NewMockPostQueryService(ctrl)
	queryService.EXPECT().FindAll(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	uc := newTestUseCase(t, queryService)
	output, err := uc.Execute(context.Background(), post.ListPostsInput{Limit: 20, Offset: -1})

	require.Error(t, err)
	assert.Nil(t, output)

	var voErr vo.Error
	require.ErrorAs(t, err, &voErr)
	assert.Equal(t, vo.ValidationErrorCode, voErr.Code())
}

func TestListPostsUseCase_LargeOffset(t *testing.T) {
	ctrl := gomock.NewController(t)

	queryService := mock_query.NewMockPostQueryService(ctrl)
	queryService.EXPECT().FindAll(gomock.Any(), 20, 999999).Return([]post.PostDto{}, 3, nil).Times(1)

	uc := newTestUseCase(t, queryService)
	output, err := uc.Execute(context.Background(), post.ListPostsInput{Limit: 20, Offset: 999999})

	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Equal(t, 3, output.Total)
	assert.Empty(t, output.Posts)
}

func TestListPostsUseCase_QueryServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)

	queryService := mock_query.NewMockPostQueryService(ctrl)
	queryService.EXPECT().FindAll(gomock.Any(), 20, 0).Return(nil, 0, errors.New("db error")).Times(1)

	uc := newTestUseCase(t, queryService)
	output, err := uc.Execute(context.Background(), post.ListPostsInput{Limit: 20, Offset: 0})

	require.Error(t, err)
	assert.Nil(t, output)
}
