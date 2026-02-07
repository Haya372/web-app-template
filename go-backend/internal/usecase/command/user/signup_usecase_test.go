package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Haya372/web-app-template/go-backend/internal/usecase/command/user"
	mock_entity "github.com/Haya372/web-app-template/go-backend/test/mock/domain/entity"
	mock_repository "github.com/Haya372/web-app-template/go-backend/test/mock/domain/entity/repository"
	mock_shared "github.com/Haya372/web-app-template/go-backend/test/mock/usecase/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestSignupUseCase_HappyCase(t *testing.T) {
	tests := []struct {
		name  string
		input user.SignupInput
	}{
		{
			name: "Success signup",
			input: user.SignupInput{
				Name:     "test",
				Email:    "test@example.com",
				Password: "password",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()

			txManager := mock_shared.NewMockTransactionManager(nil)

			userRepository := mock_repository.NewMockUserRepository(ctrl)
			mockUser := mock_entity.NewMockUser(ctrl)
			userRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(mockUser, nil).Times(1)

			usecase := user.NewSignupUseCase(userRepository, txManager)

			output, err := usecase.Execute(ctx, tt.input)

			require.NoError(t, err)
			assert.Equal(t, output.Name, tt.input.Name)
			assert.Equal(t, output.Email, tt.input.Email)
		})
	}
}

func TestSignupUseCase_FailureCase(t *testing.T) {
	tests := []struct {
		name      string
		input     user.SignupInput
		createErr error
		txError   error
	}{
		{
			name: "failed to create user",
			input: user.SignupInput{
				Name:     "test",
				Email:    "test@example.com",
				Password: "passwor",
			},
			createErr: nil,
			txError:   nil,
		},
		{
			name: "failed to save user",
			input: user.SignupInput{
				Name:     "test",
				Email:    "test@example.com",
				Password: "password",
			},
			createErr: errors.New("test"),
			txError:   nil,
		},
		{
			name: "transaction error",
			input: user.SignupInput{
				Name:     "test",
				Email:    "test@example.com",
				Password: "password",
			},
			createErr: nil,
			txError:   errors.New("test"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()

			txManager := mock_shared.NewMockTransactionManager(tt.txError)

			userRepository := mock_repository.NewMockUserRepository(ctrl)
			mockUser := mock_entity.NewMockUser(ctrl)
			userRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(mockUser, tt.createErr).AnyTimes()

			usecase := user.NewSignupUseCase(userRepository, txManager)

			output, err := usecase.Execute(ctx, tt.input)

			require.Error(t, err)
			assert.Nil(t, output)
		})
	}
}
