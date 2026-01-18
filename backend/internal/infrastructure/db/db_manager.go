package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/Haya372/go-template/backend/internal/infrastructure/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DbManager interface {
	QueriesFunc(ctx context.Context, fn func(ctx context.Context, queries sqlc.Queries) error) error
	PoolFunc(ctx context.Context, fn func(ctx context.Context, conn *pgxpool.Conn) error) error
}

type dbManagerImpl struct {
	pool *pgxpool.Pool
}

func (m *dbManagerImpl) QueriesFunc(ctx context.Context, fn func(ctx context.Context, queries sqlc.Queries) error) error {
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
		return errors.New("illegal tx value")
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

func (info *DbInfo) GetDatabaseUrl() string {
	if len(info.Dsn) > 0 {
		return info.Dsn
	}
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", info.User, info.Password, info.Host, info.Port, info.Database)
}

func NewDbManager(ctx context.Context, info DbInfo) (DbManager, error) {
	dbpool, err := pgxpool.New(ctx, info.GetDatabaseUrl())
	if err != nil {
		return nil, err
	}
	return &dbManagerImpl{
		pool: dbpool,
	}, nil
}
