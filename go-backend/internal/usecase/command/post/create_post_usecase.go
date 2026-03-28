package post

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

type CreatePostUseCase interface {
	Execute(ctx context.Context, input CreatePostInput) (*CreatePostOutput, error)
}

type CreatePostInput struct {
	UserID  uuid.UUID
	Content string
}

type CreatePostOutput struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Content   string
	CreatedAt time.Time
}

type createPostUseCaseImpl struct {
	tracer         trace.Tracer
	logger         common.Logger
	postRepository repository.PostRepository
	txManager      shared.TransactionManager
}

func (uc *createPostUseCaseImpl) Execute(ctx context.Context, input CreatePostInput) (*CreatePostOutput, error) {
	ctx, span := uc.tracer.Start(ctx, "execute")
	defer span.End()

	post, err := entity.NewPost(input.UserID, input.Content, time.Now())
	if err != nil {
		uc.logger.Error(ctx, "failed to create Post", "error", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	var created entity.Post

	err = uc.txManager.Do(ctx, func(ctx context.Context) error {
		var repoErr error

		created, repoErr = uc.postRepository.Create(ctx, post)
		if repoErr != nil {
			uc.logger.Error(ctx, "failed to save Post", "error", repoErr)

			return repoErr
		}

		return nil
	})
	if err != nil {
		uc.logger.Error(ctx, "transaction error", "error", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	return &CreatePostOutput{
		ID:        created.ID(),
		UserID:    created.UserID(),
		Content:   created.Content(),
		CreatedAt: created.CreatedAt(),
	}, nil
}

func NewCreatePostUseCase(
	postRepository repository.PostRepository,
	txManager shared.TransactionManager,
) CreatePostUseCase {
	return &createPostUseCaseImpl{
		tracer:         otel.Tracer("CreatePostUseCase"),
		logger:         common.NewLogger(),
		postRepository: postRepository,
		txManager:      txManager,
	}
}
