package integration

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDbProps struct {
	User      string
	Password  string
	Database  string
	DbDirPath string
	Schema    string
}

type TestDb interface {
	DbManager() db.DbManager
	Cleanup() error
	Terminate() error
	Pool() *pgxpool.Pool
}

type baseTestDb struct {
	pool    *pgxpool.Pool
	manager db.DbManager
	schema  string
}

func (b *baseTestDb) DbManager() db.DbManager { return b.manager }

func (b *baseTestDb) Pool() *pgxpool.Pool { return b.pool }

func (b *baseTestDb) Cleanup() error {
	return b.manager.PoolFunc(context.Background(), func(ctx context.Context, conn *pgxpool.Conn) error {
		// Truncate in dependency order: posts and user_roles reference users.
		_, err := conn.Exec(ctx, "truncate table posts, user_roles, users")

		return err
	})
}

type localTestDb struct {
	baseTestDb

	container testcontainers.Container
}

type ciTestDb struct {
	baseTestDb
}

const wailOccurrence = 2

var errEmptySchema = errors.New("TestDbProps.Schema must not be empty")

func NewTestServer(e *echo.Echo) *httptest.Server {
	return httptest.NewServer(e)
}

func NewTestDb(props TestDbProps) (TestDb, error) {
	if props.Schema == "" {
		return nil, errEmptySchema
	}

	ctx := context.Background()

	if os.Getenv("CI") == "true" {
		return newCITestDb(ctx, props)
	}

	return newLocalTestDb(ctx, props)
}

func newCITestDb(ctx context.Context, props TestDbProps) (TestDb, error) {
	pool, err := createSchemaPool(ctx, os.Getenv("DATABASE_DSN"), props.Schema)
	if err != nil {
		return nil, err
	}

	manager := db.NewDbManager(pool)

	if err = runMigrations(ctx, manager, props.DbDirPath); err != nil {
		pool.Close()

		return nil, err
	}

	return &ciTestDb{baseTestDb: baseTestDb{manager: manager, pool: pool, schema: props.Schema}}, nil
}

func newLocalTestDb(ctx context.Context, props TestDbProps) (TestDb, error) {
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
		_ = container.Terminate(ctx)

		return nil, err
	}

	pool, err := createSchemaPool(ctx, dsn, props.Schema)
	if err != nil {
		_ = container.Terminate(ctx)

		return nil, err
	}

	manager := db.NewDbManager(pool)

	if err = runMigrations(ctx, manager, props.DbDirPath); err != nil {
		pool.Close()

		_ = container.Terminate(ctx)

		return nil, err
	}

	return &localTestDb{
		baseTestDb: baseTestDb{manager: manager, pool: pool, schema: props.Schema},
		container:  container,
	}, nil
}

// createSchemaPool creates the given schema if it does not exist, then returns
// a pool whose connections have search_path set to that schema.
func createSchemaPool(ctx context.Context, dsn, schema string) (*pgxpool.Pool, error) {
	tmpPool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	_, err = tmpPool.Exec(ctx, "CREATE SCHEMA IF NOT EXISTS "+pgx.Identifier{schema}.Sanitize())
	tmpPool.Close()

	if err != nil {
		return nil, fmt.Errorf("create schema %s: %w", schema, err)
	}

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	config.ConnConfig.RuntimeParams["search_path"] = schema

	return pgxpool.NewWithConfig(ctx, config)
}

func runMigrations(ctx context.Context, manager db.DbManager, dbDirPath string) error {
	return manager.PoolFunc(ctx, func(ctx context.Context, conn *pgxpool.Conn) error {
		if err := runSQLDir(ctx, conn, path.Join(dbDirPath, "schema")); err != nil {
			return err
		}

		return runSQLDir(ctx, conn, path.Join(dbDirPath, "seeds", "master"))
	})
}


func (d *localTestDb) Terminate() error {
	d.pool.Close()

	return d.container.Terminate(context.Background())
}

func (d *ciTestDb) Terminate() error {
	ctx := context.Background()

	err := d.manager.PoolFunc(ctx, func(ctx context.Context, conn *pgxpool.Conn) error {
		_, err := conn.Exec(ctx, "DROP SCHEMA IF EXISTS "+pgx.Identifier{d.schema}.Sanitize()+" CASCADE")

		return err
	})
	if err != nil {
		slog.Warn("failed to drop schema", "schema", d.schema, "error", err)
	}

	d.pool.Close()

	return nil
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
