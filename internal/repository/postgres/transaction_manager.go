package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"service-order-avito/internal/domain/errors/repository"
)

type transactionManagerPostgres struct {
	pool *pgxpool.Pool
}

func NewTransactionManagerPostgres(pool *pgxpool.Pool) *transactionManagerPostgres {
	return &transactionManagerPostgres{pool: pool}
}

func (tm *transactionManagerPostgres) Begin(parent context.Context, fn func(ctx context.Context) error) error {
	tx, err := tm.pool.Begin(parent)
	if err != nil {
		return repository.ErrInternalError
	}

	ctx := context.WithValue(parent, txContextKey{}, tx)

	err = fn(ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

type txContextKey struct{}

func GetTx(ctx context.Context) pgx.Tx {
	tx, _ := ctx.Value(txContextKey{}).(pgx.Tx)
	return tx
}
