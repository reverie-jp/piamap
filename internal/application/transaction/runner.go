package transaction

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/reverie-jp/piamap/internal/gen/sqlc"
)

type Runner interface {
	WithTx(ctx context.Context, fn func(q sqlc.Querier) error) error
}

type RunnerImpl struct {
	pool *pgxpool.Pool
}

func NewRunner(pool *pgxpool.Pool) Runner {
	return &RunnerImpl{pool: pool}
}

func (r *RunnerImpl) WithTx(ctx context.Context, fn func(q sqlc.Querier) error) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	if err := fn(sqlc.New(tx)); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return rbErr
		}
		return err
	}
	return tx.Commit(ctx)
}
