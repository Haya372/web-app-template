package db

import (
	"context"
	"errors"

	"github.com/Haya372/go-template/backend/internal/infrastructure/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DbManager interface {
	QueriesFunc(ctx context.Context, fn func(ctx context.Context, queries sqlc.Queries) error) error
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

func NewDbManager(pool *pgxpool.Pool) DbManager {
	return &dbManagerImpl{
		pool: pool,
	}
}
