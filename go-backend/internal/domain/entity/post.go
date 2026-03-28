//go:generate mockgen -source=post.go -destination=../../../test/mock/domain/entity/mock_post.go

package entity

import (
	"time"

	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	"github.com/google/uuid"
)

// Post represents a post domain entity.
type Post interface {
	ID() uuid.UUID
	UserID() uuid.UUID
	Content() string
	CreatedAt() time.Time
}

type postImpl struct {
	id        uuid.UUID
	userId    uuid.UUID
	content   string
	createdAt time.Time
}

func (p *postImpl) ID() uuid.UUID {
	return p.id
}

func (p *postImpl) UserID() uuid.UUID {
	return p.userId
}

func (p *postImpl) Content() string {
	return p.content
}

func (p *postImpl) CreatedAt() time.Time {
	return p.createdAt
}

// NewPost creates a new Post with a generated UUID, validating the content.
func NewPost(userID uuid.UUID, content string, createdAt time.Time) (Post, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	c, err := vo.NewContent(content)
	if err != nil {
		return nil, err
	}

	return &postImpl{
		id:        id,
		userId:    userID,
		content:   string(*c),
		createdAt: createdAt,
	}, nil
}

// ReconstructPost rebuilds a Post from persisted values without validation.
func ReconstructPost(id, userID uuid.UUID, content string, createdAt time.Time) Post {
	return &postImpl{
		id:        id,
		userId:    userID,
		content:   content,
		createdAt: createdAt,
	}
}
