//go:generate mockgen -source=jwt_service.go -destination=../../../test/mock/usecase/service/mock_jwt_service.go

package service

import (
	"context"
	"time"

	"github.com/Haya372/web-app-template/go-backend/internal/domain/entity"
)

type UserAccessToken struct {
	Value     string
	ExpiresAt time.Time
}

type JwtService interface {
	GenerateUserAccessToken(ctx context.Context, user entity.User) (*UserAccessToken, error)
}
