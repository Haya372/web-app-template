package integration

import (
	"context"
	"fmt"
	"log/slog"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"sort"
	"testing"

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

type localTestDb struct {
	pool      *pgxpool.Pool
	manager   db.DbManager
	container testcontainers.Container
	schema    string
}

type ciTestDb struct {
	pool    *pgxpool.Pool
	manager db.DbManager
	schema  string
}

const wailOccurrence = 2

func NewTestServer(e *echo.Echo) *httptest.Server {
	return httptest.NewServer(e)
}

func NewTestDb(props TestDbProps) (TestDb, error) {
	ctx := context.Background()

	if props.Schema == "" {
		return nil, fmt.Errorf("TestDbProps.Schema must not be empty")
	}

	if os.Getenv("CI") == "true" {
		dsn := os.Getenv("DATABASE_DSN")

		// Step 1: create the schema using a temporary connection.
		tmpPool, err := pgxpool.New(ctx, dsn)
		if err != nil {
			return nil, err
		}
		if _, err = tmpPool.Exec(ctx, fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", pgx.Identifier{props.Schema}.Sanitize())); err != nil {
			tmpPool.Close()
			return nil, fmt.Errorf("create schema %s: %w", props.Schema, err)
		}
		tmpPool.Close()

		// Step 2: create pool with search_path set to the schema.
		config, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			return nil, err
		}
		config.ConnConfig.RuntimeParams["search_path"] = props.Schema

		pool, err := pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			return nil, err
		}

		manager := db.NewDbManager(pool)

		err = manager.PoolFunc(ctx, func(ctx context.Context, conn *pgxpool.Conn) error {
			// running migration
			if err := runSQLDir(ctx, conn, path.Join(props.DbDirPath, "schema")); err != nil {
				return err
			}

			// running seed generation
			return runSQLDir(ctx, conn, path.Join(props.DbDirPath, "seeds", "master"))
		})
		if err != nil {
			pool.Close()
			return nil, err
		}

		return &ciTestDb{manager: manager, pool: pool, schema: props.Schema}, nil
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

	// Step 1: create the schema using a temporary connection.
	tmpPool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	if _, err = tmpPool.Exec(ctx, fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", pgx.Identifier{props.Schema}.Sanitize())); err != nil {
		tmpPool.Close()
		return nil, fmt.Errorf("create schema %s: %w", props.Schema, err)
	}
	tmpPool.Close()

	// Step 2: create pool with search_path set to the schema.
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	config.ConnConfig.RuntimeParams["search_path"] = props.Schema

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	manager := db.NewDbManager(pool)

	err = manager.PoolFunc(ctx, func(ctx context.Context, conn *pgxpool.Conn) error {
		// running migration
		if err := runSQLDir(ctx, conn, path.Join(props.DbDirPath, "schema")); err != nil {
			return err
		}

		// running seed generation
		return runSQLDir(ctx, conn, path.Join(props.DbDirPath, "seeds", "master"))
	})
	if err != nil {
		pool.Close()
		_ = container.Terminate(ctx)
		return nil, err
	}

	return &localTestDb{
		manager:   manager,
		container: container,
		pool:      pool,
		schema:    props.Schema,
	}, nil
}

// WithTx begins a transaction on the given TestDb's pool and registers a
// deferred rollback via t.Cleanup. The returned context carries the
// transaction so that DbManager.QueriesFunc uses it instead of acquiring a
// new connection. Use this in repository tests to achieve per-test isolation
// without TRUNCATE.
func WithTx(t *testing.T, testDb TestDb) context.Context {
	t.Helper()

	tx, err := testDb.Pool().Begin(context.Background())
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}

	t.Cleanup(func() {
		_ = tx.Rollback(context.Background())
	})

	return db.WithTx(context.Background(), tx)
}

func (d *localTestDb) DbManager() db.DbManager {
	return d.manager
}

func (d *localTestDb) Cleanup() error {
	return d.manager.PoolFunc(context.Background(), func(ctx context.Context, conn *pgxpool.Conn) error {
		_, err := conn.Exec(ctx, "truncate table users")

		return err
	})
}

func (d *localTestDb) Terminate() error {
	d.pool.Close()

	return d.container.Terminate(context.Background())
}

func (d *localTestDb) Pool() *pgxpool.Pool {
	return d.pool
}

func (d *ciTestDb) DbManager() db.DbManager {
	return d.manager
}

func (d *ciTestDb) Cleanup() error {
	return d.manager.PoolFunc(context.Background(), func(ctx context.Context, conn *pgxpool.Conn) error {
		_, err := conn.Exec(ctx, "truncate table users")

		return err
	})
}

func (d *ciTestDb) Terminate() error {
	ctx := context.Background()

	err := d.manager.PoolFunc(ctx, func(ctx context.Context, conn *pgxpool.Conn) error {
		_, err := conn.Exec(ctx, fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE", pgx.Identifier{d.schema}.Sanitize()))

		return err
	})
	if err != nil {
		slog.Warn("failed to drop schema", "schema", d.schema, "error", err)
	}

	d.pool.Close()

	return nil
}

func (d *ciTestDb) Pool() *pgxpool.Pool {
	return d.pool
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
