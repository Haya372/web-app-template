//go:generate mockgen -source=post.go -destination=../../../test/mock/domain/entity/mock_post.go

package entity

import (
	"time"

	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	"github.com/google/uuid"
)

// Post represents a post domain entity.
type Post interface {
	Id() uuid.UUID
	UserId() uuid.UUID
	Content() string
	CreatedAt() time.Time
}

type postImpl struct {
	id        uuid.UUID
	userId    uuid.UUID
	content   string
	createdAt time.Time
}

func (p *postImpl) Id() uuid.UUID {
	return p.id
}

func (p *postImpl) UserId() uuid.UUID {
	return p.userId
}

func (p *postImpl) Content() string {
	return p.content
}

func (p *postImpl) CreatedAt() time.Time {
	return p.createdAt
}

// NewPost creates a new Post with a generated UUID, validating the content.
func NewPost(userId uuid.UUID, content string, createdAt time.Time) (Post, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	c, err := vo.NewContent(content)
	if err != nil {
		return nil, err
	}

	return &postImpl{
		id:        id,
		userId:    userId,
		content:   string(*c),
		createdAt: createdAt,
	}, nil
}

// ReconstructPost rebuilds a Post from persisted values without validation.
func ReconstructPost(id, userId uuid.UUID, content string, createdAt time.Time) Post {
	return &postImpl{
		id:        id,
		userId:    userId,
		content:   content,
		createdAt: createdAt,
	}
}
