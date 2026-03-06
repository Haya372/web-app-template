//go:generate mockgen -source=list_users_query.go -destination=../../../../test/mock/usecase/query/mock_user_query_service.go

package user

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// UserDto is a read-only projection of a user; password_hash is intentionally excluded.
type UserDto struct {
	Id        uuid.UUID
	Name      string
	Email     string
	Status    string
	CreatedAt time.Time
}

// UserQueryService is the port for fetching user projections from the data store.
type UserQueryService interface {
	FindAll(ctx context.Context, limit, offset int) ([]UserDto, int, error)
}

// ListUsersInput holds the validated parameters for the list-users query.
type ListUsersInput struct {
	Limit  int
	Offset int
}

// ListUsersOutput is the result returned by ListUsersUseCase.
type ListUsersOutput struct {
	Users []UserDto
	Total int
}

// ListUsersUseCase is the application use case for retrieving a paginated user list.
type ListUsersUseCase interface {
	Execute(ctx context.Context, input ListUsersInput) (*ListUsersOutput, error)
}
