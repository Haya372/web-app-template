package common_mock

import (
	"context"

	"github.com/Haya372/web-app-template/backend/internal/common"
)

type mockTransactionManager struct {
	err error
}

func (tx *mockTransactionManager) Do(ctx context.Context, f func(ctx context.Context) error) error {
	if tx.err != nil {
		return tx.err
	}
	return f(ctx)
}

func NewMockTransactionManager(err error) common.TransactionManager {
	return &mockTransactionManager{
		err: err,
	}
}
