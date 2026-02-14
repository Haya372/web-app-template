package user_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Haya372/web-app-template/go-backend/internal/domain/entity/repository"
	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	"github.com/Haya372/web-app-template/go-backend/internal/usecase/command/user"
	"github.com/Haya372/web-app-template/go-backend/internal/usecase/service"
	mock_entity "github.com/Haya372/web-app-template/go-backend/test/mock/domain/entity"
	mock_repository "github.com/Haya372/web-app-template/go-backend/test/mock/domain/entity/repository"
	mock_service "github.com/Haya372/web-app-template/go-backend/test/mock/usecase/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestLoginUseCase_HappyCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	userRepository := mock_repository.NewMockUserRepository(ctrl)
	mockUser := mock_entity.NewMockUser(ctrl)
	userID := uuid.New()
	expiresAt := time.Date(2026, 2, 14, 12, 0, 0, 0, time.UTC)

	userRepository.EXPECT().FindByEmail(gomock.Any(), "test@example.com").Return(mockUser, nil).Times(1)
	mockUser.EXPECT().Status().Return(vo.UserStatusActive).Times(1)
	mockUser.EXPECT().ComparePassword("password").Return(true, nil).Times(1)
	mockUser.EXPECT().Id().Return(userID).Times(1)
	mockUser.EXPECT().Name().Return("Test").Times(1)
	mockUser.EXPECT().Email().Return("test@example.com").Times(1)

	tokenGenerator := mock_service.NewMockJwtService(ctrl)
	tokenGenerator.EXPECT().
		GenerateUserAccessToken(gomock.Any(), mockUser).
		Return(&service.UserAccessToken{
			Value:     "token",
			ExpiresAt: expiresAt,
		}, nil).
		Times(1)

	usecase := user.NewLoginUseCase(userRepository, tokenGenerator)

	output, err := usecase.Execute(ctx, user.LoginInput{
		Email:    "test@example.com",
		Password: "password",
	})

	require.NoError(t, err)
	assert.Equal(t, "token", output.Token)
	assert.Equal(t, expiresAt, output.ExpiresAt)
	assert.Equal(t, userID.String(), output.UserID)
	assert.Equal(t, "Test", output.UserName)
	assert.Equal(t, "test@example.com", output.UserEmail)
}

func TestLoginUseCase_FailureCase(t *testing.T) {
	tests := []struct {
		name        string
		setupMocks  func(ctrl *gomock.Controller) (*mock_repository.MockUserRepository, *mock_service.MockJwtService)
		input       user.LoginInput
		assertError func(t *testing.T, err error)
	}{
		{
			name: "user not found",
			setupMocks: func(ctrl *gomock.Controller) (*mock_repository.MockUserRepository, *mock_service.MockJwtService) {
				userRepository := mock_repository.NewMockUserRepository(ctrl)
				userRepository.EXPECT().
					FindByEmail(gomock.Any(), "missing@example.com").
					Return(nil, repository.ErrUserNotFound).
					Times(1)

				return userRepository, mock_service.NewMockJwtService(ctrl)
			},
			input: user.LoginInput{
				Email:    "missing@example.com",
				Password: "password",
			},
			assertError: assertUnauthorizedError,
		},
		{
			name: "user not active",
			setupMocks: func(ctrl *gomock.Controller) (*mock_repository.MockUserRepository, *mock_service.MockJwtService) {
				userRepository := mock_repository.NewMockUserRepository(ctrl)
				mockUser := mock_entity.NewMockUser(ctrl)
				userRepository.EXPECT().FindByEmail(gomock.Any(), "test@example.com").Return(mockUser, nil).Times(1)
				mockUser.EXPECT().Status().Return(vo.UserStatusDeleted).Times(1)

				return userRepository, mock_service.NewMockJwtService(ctrl)
			},
			input: user.LoginInput{
				Email:    "test@example.com",
				Password: "password",
			},
			assertError: assertUnauthorizedError,
		},
		{
			name: "password mismatch",
			setupMocks: func(ctrl *gomock.Controller) (*mock_repository.MockUserRepository, *mock_service.MockJwtService) {
				userRepository := mock_repository.NewMockUserRepository(ctrl)
				mockUser := mock_entity.NewMockUser(ctrl)
				userRepository.EXPECT().FindByEmail(gomock.Any(), "test@example.com").Return(mockUser, nil).Times(1)
				mockUser.EXPECT().Status().Return(vo.UserStatusActive).Times(1)
				mockUser.EXPECT().ComparePassword("wrong").Return(false, nil).Times(1)

				return userRepository, mock_service.NewMockJwtService(ctrl)
			},
			input: user.LoginInput{
				Email:    "test@example.com",
				Password: "wrong",
			},
			assertError: assertUnauthorizedError,
		},
		{
			name: "token generation error",
			setupMocks: func(ctrl *gomock.Controller) (*mock_repository.MockUserRepository, *mock_service.MockJwtService) {
				userRepository := mock_repository.NewMockUserRepository(ctrl)
				mockUser := mock_entity.NewMockUser(ctrl)
				userRepository.EXPECT().FindByEmail(gomock.Any(), "test@example.com").Return(mockUser, nil).Times(1)
				mockUser.EXPECT().Status().Return(vo.UserStatusActive).Times(1)
				mockUser.EXPECT().ComparePassword("password").Return(true, nil).Times(1)

				jwtService := mock_service.NewMockJwtService(ctrl)
				jwtService.EXPECT().
					GenerateUserAccessToken(gomock.Any(), mockUser).
					Return(nil, errors.New("token error")).
					Times(1)

				return userRepository, jwtService
			},
			input: user.LoginInput{
				Email:    "test@example.com",
				Password: "password",
			},
			assertError: func(t *testing.T, err error) {
				t.Helper()
				require.Error(t, err)

				var baseErr vo.Error
				assert.NotErrorAs(t, err, &baseErr)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()

			userRepository, tokenGenerator := tt.setupMocks(ctrl)
			usecase := user.NewLoginUseCase(userRepository, tokenGenerator)

			output, err := usecase.Execute(ctx, tt.input)

			require.Error(t, err)
			assert.Nil(t, output)
			tt.assertError(t, err)
		})
	}
}

func assertUnauthorizedError(t *testing.T, err error) {
	t.Helper()

	var baseErr vo.Error
	require.ErrorAs(t, err, &baseErr)
	assert.Equal(t, vo.InvalidCredentialErrorCode, baseErr.Code())
}
