//go:build wireinject

package integration

import (
	"context"
	"net/http/httptest"

	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/db"
	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/http"
	infraquery "github.com/Haya372/web-app-template/go-backend/internal/infrastructure/query"
	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/repository"
	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/service"
	commandpost "github.com/Haya372/web-app-template/go-backend/internal/usecase/command/post"
	"github.com/Haya372/web-app-template/go-backend/internal/usecase/command/user"
	querypost "github.com/Haya372/web-app-template/go-backend/internal/usecase/query/post"
	queryuser "github.com/Haya372/web-app-template/go-backend/internal/usecase/query/user"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
)

var repositorySet = wire.NewSet(
	repository.NewUserRepository,
	repository.NewPostRepository,
)

var authSet = wire.NewSet(
	service.NewJwtService,
)

var usecaseSet = wire.NewSet(
	user.NewSignupUseCase,
	user.NewLoginUseCase,
	commandpost.NewCreatePostUseCase,
)

var querySet = wire.NewSet(
	infraquery.NewUserQueryService,
	infraquery.NewPostQueryService,
	repository.NewUserPermissionRepository,
	queryuser.NewListUsersUseCase,
	querypost.NewListPostsUseCase,
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
		authSet,
		usecaseSet,
		querySet,
		dbSet,
		httpSet,
		testServerSet,
	)

	return nil, nil
}
