package common

import "context"

type TransactionManager interface {
	Do(ctx context.Context, f func(ctx context.Context) error) error
}
