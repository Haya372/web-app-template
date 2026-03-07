//go:generate mockgen -source=post_repository.go -destination=../../../../test/mock/domain/entity/repository/mock_post_repository.go

package repository

import (
	"context"

	"github.com/Haya372/web-app-template/go-backend/internal/domain/entity"
)

type PostRepository interface {
	Create(ctx context.Context, post entity.Post) (entity.Post, error)
}
