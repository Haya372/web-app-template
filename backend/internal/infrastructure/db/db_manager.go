package db

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/Haya372/web-app-template/backend/internal/infrastructure/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrIllegalTx = errors.New("illegal tx value")
)

type DbManager interface {
	QueriesFunc(ctx context.Context, fn func(ctx context.Context, queries sqlc.Queries) error) error
	PoolFunc(ctx context.Context, fn func(ctx context.Context, conn *pgxpool.Conn) error) error
}

type dbManagerImpl struct {
	pool *pgxpool.Pool
}

func (m *dbManagerImpl) QueriesFunc(
	ctx context.Context,
	fn func(ctx context.Context, queries sqlc.Queries) error) error {
	val := ctx.Value(txKey)
	if val == nil {
		conn, err := m.pool.Acquire(ctx)
		if err != nil {
			return err
		}
		defer conn.Release()

		queries := sqlc.New(conn)

		return fn(ctx, *queries)
	}

	tx, ok := val.(pgx.Tx)
	if !ok {
		return ErrIllegalTx
	}

	queries := sqlc.New(tx)

	return fn(ctx, *queries)
}

func (m *dbManagerImpl) PoolFunc(ctx context.Context, fn func(ctx context.Context, conn *pgxpool.Conn) error) error {
	conn, err := m.pool.Acquire(ctx)
	if err != nil {
		return err
	}

	defer conn.Release()

	return fn(ctx, conn)
}

type DbInfo struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Dsn      string
}

func NewDbInfo() DbInfo {
	dsn := os.Getenv("DATABASE_DSN")

	return DbInfo{
		Dsn: dsn,
	}
}

func (info *DbInfo) GetDatabaseUrl() string {
	if len(info.Dsn) > 0 {
		return info.Dsn
	}

	return fmt.Sprintf(
		"postgresql://%s:%s@%s/%s",
		info.User, info.Password, net.JoinHostPort(info.Host, info.Port), info.Database)
}

func NewDbManager(pool *pgxpool.Pool) DbManager {
	return &dbManagerImpl{
		pool: pool,
	}
}
