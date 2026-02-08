//go:build wireinject

package integration

import (
	"context"
	"net/http/httptest"

	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/db"
	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/http"
	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/repository"
	"github.com/Haya372/web-app-template/go-backend/internal/usecase/command/user"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
)

var repositorySet = wire.NewSet(
	repository.NewUserRepository,
)

var usecaseSet = wire.NewSet(
	user.NewSignupUseCase,
)

var dbSet = wire.NewSet(
	db.NewDbManager,
	db.NewTransactionManger,
)

var httpSet = wire.NewSet(
	http.NewRouter,
	http.NewServer,
)

var testServerSet = wire.NewSet(
	NewTestServer,
)

func InitializeTestServer(ctx context.Context, pool *pgxpool.Pool) (*httptest.Server, error) {
	wire.Build(
		repositorySet,
		usecaseSet,
		dbSet,
		httpSet,
		testServerSet,
	)
	return nil, nil
}
