//go:build wireinject

package di

import (
	"context"

	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/db"
	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/http"
	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/repository"
	"github.com/Haya372/web-app-template/go-backend/internal/usecase/command/user"
	"github.com/google/wire"
)

var repositorySet = wire.NewSet(
	repository.NewUserRepository,
)

var usecaseSet = wire.NewSet(
	user.NewSignupUseCase,
)

var dbSet = wire.NewSet(
	db.NewDbPool,
	db.NewDbInfo,
	db.NewDbManager,
	db.NewTransactionManger,
)

var httpSet = wire.NewSet(
	wire.Struct(new(http.Router), "*"),
	http.NewEchoConfig,
	http.NewServer,
	wire.Struct(new(http.Server), "*"),
)

func InitializeServer(ctx context.Context) (*http.Server, error) {
	wire.Build(
		repositorySet,
		usecaseSet,
		dbSet,
		httpSet,
	)
	return nil, nil
}
