package post_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Haya372/web-app-template/go-backend/internal/usecase/command/post"
	mock_entity "github.com/Haya372/web-app-template/go-backend/test/mock/domain/entity"
	mock_repository "github.com/Haya372/web-app-template/go-backend/test/mock/domain/entity/repository"
	mock_shared "github.com/Haya372/web-app-template/go-backend/test/mock/usecase/shared"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreatePostUseCase_HappyCase(t *testing.T) {
	userId := uuid.New()
	postId := uuid.New()
	createdAt := time.Date(2026, 3, 8, 0, 0, 0, 0, time.UTC)
	content := "Hello, world!"

	tests := []struct {
		name  string
		input post.CreatePostInput
	}{
		{
			name: "Success to create post",
			input: post.CreatePostInput{
				UserId:  userId,
				Content: content,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()

			txManager := mock_shared.NewMockTransactionManager(nil)
			postRepository := mock_repository.NewMockPostRepository(ctrl)
			mockPost := mock_entity.NewMockPost(ctrl)

			mockPost.EXPECT().Id().Return(postId).AnyTimes()
			mockPost.EXPECT().UserId().Return(userId).AnyTimes()
			mockPost.EXPECT().Content().Return(content).AnyTimes()
			mockPost.EXPECT().CreatedAt().Return(createdAt).AnyTimes()

			postRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(mockPost, nil).Times(1)

			usecase := post.NewCreatePostUseCase(postRepository, txManager)
			output, err := usecase.Execute(ctx, tt.input)

			require.NoError(t, err)
			require.NotNil(t, output)
			assert.Equal(t, postId, output.Id)
			assert.Equal(t, userId, output.UserId)
			assert.Equal(t, content, output.Content)
			assert.Equal(t, createdAt, output.CreatedAt)
		})
	}
}

func TestCreatePostUseCase_FailureCase(t *testing.T) {
	userId := uuid.New()

	tests := []struct {
		name      string
		input     post.CreatePostInput
		repoErr   error
		txErr     error
	}{
		{
			name: "empty content returns error",
			input: post.CreatePostInput{
				UserId:  userId,
				Content: "",
			},
		},
		{
			name: "repository error propagates",
			input: post.CreatePostInput{
				UserId:  userId,
				Content: "valid content",
			},
			repoErr: errors.New("db error"),
		},
		{
			name: "transaction error propagates",
			input: post.CreatePostInput{
				UserId:  userId,
				Content: "valid content",
			},
			txErr: errors.New("tx error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()

			txManager := mock_shared.NewMockTransactionManager(tt.txErr)
			postRepository := mock_repository.NewMockPostRepository(ctrl)
			mockPost := mock_entity.NewMockPost(ctrl)

			postRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(mockPost, tt.repoErr).AnyTimes()

			usecase := post.NewCreatePostUseCase(postRepository, txManager)
			output, err := usecase.Execute(ctx, tt.input)

			require.Error(t, err)
			assert.Nil(t, output)
		})
	}
}
