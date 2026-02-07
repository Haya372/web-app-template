package user

import (
	"context"
	"time"

	"github.com/Haya372/web-app-template/go-backend/internal/common"
	"github.com/Haya372/web-app-template/go-backend/internal/domain/entity"
	"github.com/Haya372/web-app-template/go-backend/internal/domain/entity/repository"
	"github.com/Haya372/web-app-template/go-backend/internal/usecase/shared"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
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
	tracer         trace.Tracer
	logger         common.Logger
	userRepository repository.UserRepository
	txManager      shared.TransactionManager
}

func (uc *signupUseCaseImpl) Execute(ctx context.Context, input SignupInput) (*SignupOutput, error) {
	ctx, span := uc.tracer.Start(ctx, "execute")
	defer span.End()

	user, err := entity.NewUser(input.Email, input.Password, input.Name, time.Now())
	if err != nil {
		uc.logger.Error(ctx, "failed to create User", "error", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	err = uc.txManager.Do(ctx, func(ctx context.Context) error {
		_, err := uc.userRepository.Create(ctx, user)
		if err != nil {
			uc.logger.Error(ctx, "failed to save User", "error", err)

			return err
		}

		return nil
	})
	if err != nil {
		uc.logger.Error(ctx, "transaction error", "error", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	return &SignupOutput{
		Id:        user.Id(),
		Name:      user.Name(),
		Email:     user.Email(),
		CreatedAt: user.CreatedAt(),
	}, nil
}

func NewSignupUseCase(userRepository repository.UserRepository, txManager shared.TransactionManager) SingupUseCase {
	return &signupUseCaseImpl{
		tracer:         otel.Tracer("SignupUseCase"),
		logger:         common.NewLogger(),
		userRepository: userRepository,
		txManager:      txManager,
	}
}
