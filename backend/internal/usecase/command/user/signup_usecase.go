package user

import (
	"context"
	"time"

	"github.com/Haya372/web-app-template/backend/internal/common"
	"github.com/Haya372/web-app-template/backend/internal/domain/entity"
	"github.com/Haya372/web-app-template/backend/internal/domain/entity/repository"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type SingupUseCase interface {
	Execute(ctx context.Context, input SignupInput) (*SignupOutput, error)
}

type SignupInput struct {
	Name     string
	Email    string
	Password string
}

type SignupOutput struct {
	Id        uuid.UUID
	Name      string
	Email     string
	CreatedAt time.Time
}

type signupUseCaseImpl struct {
	logger         common.Logger
	userRepository repository.UserRepository
	txManager      common.TransactionManager
}

func (uc *signupUseCaseImpl) Execute(ctx context.Context, input SignupInput) (*SignupOutput, error) {
	user, err := entity.NewUser(input.Email, input.Password, input.Name, time.Now())
	if err != nil {
		uc.logger.Error(ctx, "failed to create User", "error", err)

		return nil, errors.Wrap(err, "signupUseCaseImpl: failed to create User")
	}

	err = uc.txManager.Do(ctx, func(ctx context.Context) error {
		_, err := uc.userRepository.Create(ctx, user)
		if err != nil {
			uc.logger.Error(ctx, "failed to save User", "error", err)

			return errors.Wrap(err, "signupUseCaseImpl: failed to save user")
		}

		return nil
	})
	if err != nil {
		uc.logger.Error(ctx, "transaction error", "error", err)

		return nil, errors.Wrap(err, "signupUseCaseImpl: Transaction error")
	}

	return &SignupOutput{
		Id:        user.Id(),
		Name:      user.Name(),
		Email:     user.Email(),
		CreatedAt: user.CreatedAt(),
	}, nil
}

func NewSignupUseCase(userRepository repository.UserRepository, txManager common.TransactionManager) SingupUseCase {
	return &signupUseCaseImpl{
		logger:         common.NewLogger(),
		userRepository: userRepository,
		txManager:      txManager,
	}
}
