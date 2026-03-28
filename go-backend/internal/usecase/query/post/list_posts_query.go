//go:generate mockgen -source=list_posts_query.go -destination=../../../../test/mock/usecase/query/mock_post_query_service.go

package post

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// PostDto is a read-only projection of a post.
type PostDto struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Content   string
	CreatedAt time.Time
}

// PostQueryService is the port for fetching post projections from the data store.
type PostQueryService interface {
	// FindAll returns a paginated list of posts and the total count.
	// The returned slice is never nil; an empty table returns a zero-length slice.
	FindAll(ctx context.Context, limit, offset int) ([]PostDto, int, error)
}

// ListPostsInput holds the validated parameters for the list-posts query.
type ListPostsInput struct {
	Limit  int
	Offset int
}

// ListPostsOutput is the result returned by ListPostsUseCase.
type ListPostsOutput struct {
	Posts []PostDto
	Total int
}

// ListPostsUseCase is the application use case for retrieving a paginated post list.
type ListPostsUseCase interface {
	Execute(ctx context.Context, input ListPostsInput) (*ListPostsOutput, error)
}
