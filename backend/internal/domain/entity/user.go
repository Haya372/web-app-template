//go:generate mockgen -source=user.go -destination=../../../test/mock/domain/entity/mock_user.go

package entity

import (
	"time"

	"github.com/Haya372/web-app-template/backend/internal/domain/vo"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type User interface {
	Id() uuid.UUID
	Email() string
	PasswordHash() []byte
	Name() string
	CreatedAt() time.Time
}

type userImpl struct {
	id           uuid.UUID
	email        string
	passwordHash []byte
	name         string
	createdAt    time.Time
}

func (u *userImpl) Id() uuid.UUID {
	return u.id
}

func (u *userImpl) Email() string {
	return u.email
}

func (u *userImpl) PasswordHash() []byte {
	return u.passwordHash
}

func (u *userImpl) Name() string {
	return u.name
}

func (u *userImpl) CreatedAt() time.Time {
	return u.createdAt
}

func NewUser(email, rawPassword, name string, createdAt time.Time) (User, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, errors.Wrap(err, "NewUser: failed to generate ID")
	}

	password, err := vo.NewPassword(rawPassword)
	if err != nil {
		return nil, errors.Wrap(err, "NewUser: illegal password")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "NewUser: failed to generate password hash")
	}

	return &userImpl{
		id:           id,
		email:        email,
		passwordHash: passwordHash,
		name:         name,
		createdAt:    createdAt,
	}, nil
}

func ReconstructUser(id uuid.UUID, email string, passwordHash []byte, name string, createdAt time.Time) User {
	return &userImpl{
		id:           id,
		email:        email,
		passwordHash: passwordHash,
		name:         name,
		createdAt:    createdAt,
	}
}
