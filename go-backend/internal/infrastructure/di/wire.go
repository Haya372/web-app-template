//go:build wireinject

package di

import (
	"context"

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
	db.NewDbPool,
	db.NewDbInfo,
	db.NewDbManager,
	db.NewTransactionManger,
)

var httpSet = wire.NewSet(
	http.NewRouter,
	http.NewEchoConfig,
	http.NewServer,
	wire.Struct(new(http.Server), "*"),
)

func InitializeServer(ctx context.Context) (*http.Server, error) {
	wire.Build(
		repositorySet,
		authSet,
		usecaseSet,
		querySet,
		dbSet,
		httpSet,
	)

	return nil, nil
}
