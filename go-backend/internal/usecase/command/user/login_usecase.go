package user

import (
	"context"
	"errors"
	"time"

	"github.com/Haya372/web-app-template/go-backend/internal/common"
	"github.com/Haya372/web-app-template/go-backend/internal/domain/entity/repository"
	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	"github.com/Haya372/web-app-template/go-backend/internal/usecase/service"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type LoginUseCase interface {
	Execute(ctx context.Context, input LoginInput) (*LoginOutput, error)
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	Token     string
	ExpiresAt time.Time
	UserID    string
	UserName  string
	UserEmail string
}

type loginUseCaseImpl struct {
	tracer         trace.Tracer
	logger         common.Logger
	userRepository repository.UserRepository
	jwtService     service.JwtService
}

var (
	errUserNotActive    = errors.New("user is not active")
	errPasswordMismatch = errors.New("password mismatch")
)

func (uc *loginUseCaseImpl) Execute(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	ctx, span := uc.tracer.Start(ctx, "execute")
	defer span.End()

	user, err := uc.userRepository.FindByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, vo.NewUnauthorizedError("invalid credential", nil, err)
		}

		uc.logger.Error(ctx, "failed to find user by email", "error", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	if !user.Status().IsActive() {
		return nil, vo.NewUnauthorizedError("invalid credential", nil, errUserNotActive)
	}

	match, err := user.ComparePassword(input.Password)
	if err != nil {
		uc.logger.Error(ctx, "failed to compare password", "error", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	if !match {
		return nil, vo.NewUnauthorizedError("invalid credential", nil, errPasswordMismatch)
	}

	token, err := uc.jwtService.GenerateUserAccessToken(ctx, user)
	if err != nil {
		uc.logger.Error(ctx, "failed to generate token", "error", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	return &LoginOutput{
		Token:     token.Value,
		ExpiresAt: token.ExpiresAt,
		UserID:    user.Id().String(),
		UserName:  user.Name(),
		UserEmail: user.Email(),
	}, nil
}

func NewLoginUseCase(
	userRepository repository.UserRepository,
	jwtService service.JwtService,
) LoginUseCase {
	return &loginUseCaseImpl{
		tracer:         otel.Tracer("LoginUseCase"),
		logger:         common.NewLogger(),
		userRepository: userRepository,
		jwtService:     jwtService,
	}
}
