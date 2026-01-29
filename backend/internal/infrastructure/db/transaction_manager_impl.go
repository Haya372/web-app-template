package db

import (
	"context"

	"github.com/Haya372/web-app-template/backend/internal/common"
	"github.com/Haya372/web-app-template/backend/internal/usecase/shared"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type transactionManagerImpl struct {
	logger common.Logger
	pool   *pgxpool.Pool
}

type txKeyStruct struct{}

var txKey = txKeyStruct{}

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

	defer func() {
		err := tx.Rollback(ctx)
		if err != nil {
			txm.logger.Warn(ctx, "transaction rollback failed", "err", err)
		}
	}()

	childCtx := context.WithValue(ctx, txKey, tx)

	if err := f(childCtx); err != nil {
		txm.logger.Debug(childCtx, "rollback transaction")

		return err
	}

	return tx.Commit(ctx)
}

func NewTransactionManger(pool *pgxpool.Pool) shared.TransactionManager {
	return &transactionManagerImpl{
		logger: common.NewLogger(),
		pool:   pool,
	}
}
