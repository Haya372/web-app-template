package integration

import (
	"context"
	"fmt"
	"log/slog"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"

	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDbProps struct {
	User     string
	Password string
	Database string
	DdlPath  string
}

type TestDb interface {
	DbManager() db.DbManager
	Cleanup() error
	Terminate() error
	Pool() *pgxpool.Pool
}

type localTestDb struct {
	pool      *pgxpool.Pool
	manager   db.DbManager
	container testcontainers.Container
}

type ciTestDb struct {
	pool    *pgxpool.Pool
	manager db.DbManager
}

const wailOccurrence = 2

func NewTestServer(e *echo.Echo) *httptest.Server {
	return httptest.NewServer(e)
}

func NewTestDb(props TestDbProps) (TestDb, error) {
	ctx := context.Background()

	if os.Getenv("CI") == "true" {
		pool, err := db.NewDbPool(ctx)
		if err != nil {
			return nil, err
		}

		manager := db.NewDbManager(pool)

		err = manager.PoolFunc(ctx, func(ctx context.Context, conn *pgxpool.Conn) error {
			return runSQLDir(ctx, conn, props.DdlPath)
		})
		if err != nil {
			return nil, err
		}

		return &ciTestDb{manager: manager, pool: pool}, nil
	}

	container, err := postgres.Run(ctx,
		"postgres:18.1",
		postgres.WithDatabase(props.Database),
		postgres.WithUsername(props.User),
		postgres.WithPassword(props.Password),
		testcontainers.WithAdditionalWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(wailOccurrence),
			wait.ForListeningPort("5432/tcp"),
		),
	)
	if err != nil {
		return nil, err
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	manager := db.NewDbManager(pool)

	err = manager.PoolFunc(ctx, func(ctx context.Context, conn *pgxpool.Conn) error {
		return runSQLDir(ctx, conn, props.DdlPath)
	})
	if err != nil {
		return nil, err
	}

	return &localTestDb{
		manager:   manager,
		container: container,
		pool:      pool,
	}, nil
}

func (db *localTestDb) DbManager() db.DbManager {
	return db.manager
}

func (db *localTestDb) Cleanup() error {
	return db.manager.PoolFunc(context.Background(), func(ctx context.Context, conn *pgxpool.Conn) error {
		_, err := conn.Exec(ctx, "truncate table users")

		return err
	})
}

func (db *localTestDb) Terminate() error {
	ctx := context.Background()

	err := db.manager.PoolFunc(ctx, func(ctx context.Context, conn *pgxpool.Conn) error {
		if _, err := conn.Exec(ctx, "drop table users"); err != nil {
			return err
		}

		_, err := conn.Exec(ctx, "drop table user_statuses")

		return err
	})
	if err != nil {
		slog.Warn("failed to terminate table", "error", err)
	}

	return db.container.Terminate(context.Background())
}

func (db *localTestDb) Pool() *pgxpool.Pool {
	return db.pool
}

func (db *ciTestDb) DbManager() db.DbManager {
	return db.manager
}

func (db *ciTestDb) Cleanup() error {
	return db.manager.PoolFunc(context.Background(), func(ctx context.Context, conn *pgxpool.Conn) error {
		_, err := conn.Exec(ctx, "truncate table users")

		return err
	})
}

func (db *ciTestDb) Terminate() error {
	ctx := context.Background()

	return db.manager.PoolFunc(ctx, func(ctx context.Context, conn *pgxpool.Conn) error {
		if _, err := conn.Exec(ctx, "drop table users"); err != nil {
			return err
		}

		_, err := conn.Exec(ctx, "drop table user_statuses")

		return err
	})
}

func (db *ciTestDb) Pool() *pgxpool.Pool {
	return db.pool
}

func runSQLDir(ctx context.Context, conn *pgxpool.Conn, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var files []string

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		if filepath.Ext(e.Name()) == ".sql" {
			files = append(files, filepath.Join(dir, e.Name()))
		}
	}

	sort.Strings(files)

	for _, f := range files {
		b, err := os.ReadFile(filepath.Clean(f))
		if err != nil {
			return fmt.Errorf("read %s: %w", f, err)
		}

		if _, err := conn.Exec(ctx, string(b)); err != nil {
			return fmt.Errorf("exec %s: %w", f, err)
		}
	}

	return nil
}
