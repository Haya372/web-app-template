package db

import (
	"context"

	"github.com/Haya372/go-template/backend/internal/common"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type transactionManagerImpl struct {
	logger common.Logger
	pool   *pgxpool.Pool
}

const txKey = "transaction"

func (txm *transactionManagerImpl) Do(ctx context.Context, f func(ctx context.Context) error) error {
	txm.logger.Debug(ctx, "start transaction")
	conn, err := txm.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	childCtx := context.WithValue(ctx, txKey, tx)

	if err := f(childCtx); err != nil {
		txm.logger.Debug(childCtx, "rollback transaction")
		return err
	}

	return tx.Commit(ctx)
}

func NewTransactionManger(pool *pgxpool.Pool) common.TransactionManager {
	return &transactionManagerImpl{
		logger: common.NewLogger(),
		pool:   pool,
	}
}
