//go:generate mockgen -source=user.go -destination=../../../test/mock/domain/entity/mock_user.go

package entity

import (
	"errors"
	"time"

	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	errIllegalName = errors.New("illegal name")
)

type User interface {
	Id() uuid.UUID
	Email() string
	PasswordHash() []byte
	ComparePassword(raw string) (bool, error)
	Name() string
	CreatedAt() time.Time
	Status() vo.UserStatus
	UpdateStatus(target vo.UserStatus) (User, error)
}

type userImpl struct {
	id           uuid.UUID
	email        string
	passwordHash []byte
	name         string
	createdAt    time.Time
	status       vo.UserStatus
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

func (u *userImpl) ComparePassword(raw string) (bool, error) {
	password, err := vo.NewPassword(raw)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword(u.passwordHash, []byte(*password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (u *userImpl) Name() string {
	return u.name
}

func (u *userImpl) CreatedAt() time.Time {
	return u.createdAt
}

func (u *userImpl) Status() vo.UserStatus {
	return u.status
}

func NewUser(email, rawPassword, name string, createdAt time.Time) (User, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	if len(name) == 0 {
		return nil, vo.NewValidationError("name is required", nil, errIllegalName)
	}

	// TODO: implement email validation

	password, err := vo.NewPassword(rawPassword)
	if err != nil {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &userImpl{
		id:           id,
		email:        email,
		passwordHash: passwordHash,
		name:         name,
		createdAt:    createdAt,
		status:       vo.UserStatusActive,
	}, nil
}

func ReconstructUser(
	id uuid.UUID,
	email string,
	passwordHash []byte,
	name string,
	status vo.UserStatus,
	createdAt time.Time,
) User {
	return &userImpl{
		id:           id,
		email:        email,
		passwordHash: passwordHash,
		name:         name,
		createdAt:    createdAt,
		status:       status,
	}
}

var errInvalidUserStatusTransition = errors.New("invalid user status transition")

func (u *userImpl) UpdateStatus(target vo.UserStatus) (User, error) {
	if target == "" {
		return nil, vo.NewValidationError("status is required", nil, errInvalidUserStatusTransition)
	}

	if u.status == target {
		return nil, vo.NewValidationError("status is not changed", map[string]any{
			"status": target.String(),
		}, errInvalidUserStatusTransition)
	}

	if u.status == vo.UserStatusDeleted {
		return nil, vo.NewValidationError("deleted user status cannot transition", map[string]any{
			"current": u.status.String(),
			"target":  target.String(),
		}, errInvalidUserStatusTransition)
	}

	return &userImpl{
		id:           u.id,
		email:        u.email,
		passwordHash: u.passwordHash,
		name:         u.name,
		createdAt:    u.createdAt,
		status:       target,
	}, nil
}
