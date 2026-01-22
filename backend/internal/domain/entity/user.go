package entity

import (
	"time"

	"github.com/google/uuid"
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

func NewUser(email, password, name string, createdAt time.Time) (User, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
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
