package repositorywrapper

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kotoproger/exchange/internal/repository"
)

type RepositoryPool interface {
	GetRepository(ctx context.Context) (
		repo repository.Querier,
		commit func() error,
		rollback func(),
		release func(),
		err error,
	)
}

type Wrapper struct {
	Pool *pgxpool.Pool
	Repo *repository.Queries
}

func (w *Wrapper) GetRepository(ctx context.Context) (
	repo repository.Querier,
	commit func() error,
	rollback func(),
	release func(),
	err error,
) {
	conn, aqerror := w.Pool.Acquire(ctx)
	if aqerror != nil {
		return nil, nil, nil, nil, fmt.Errorf("get repository: %w", aqerror)
	}
	transaction, err := conn.BeginTx(
		ctx,
		pgx.TxOptions{IsoLevel: pgx.ReadCommitted},
	)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("start transaction: %w", err)
	}

	release = func(conn *pgxpool.Conn) func() {
		return func() {
			defer conn.Release()
		}
	}(conn)

	commit = func(tr pgx.Tx, ctx context.Context) func() error {
		return func() error {
			return tr.Commit(ctx)
		}
	}(transaction, ctx)

	rollback = func(tr pgx.Tx, ctx context.Context) func() {
		return func() {
			defer tr.Rollback(ctx)
		}
	}(transaction, ctx)

	return w.Repo.WithTx(transaction), commit, rollback, release, nil
}
