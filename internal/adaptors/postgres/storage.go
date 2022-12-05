package postgres

import (
	"context"

	"subscription/internal/adaptors"

	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Storage {
	return &Storage{db: pool}
}

func (s *Storage) Connection(ctx context.Context) (adaptors.ConnectionWrapper, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context, command func(ctx context.Context, q pgxtype.Querier) error) error {
		defer conn.Release()
		return command(ctx, conn)
	}, nil
}

type TxOptionFunc func(opt pgx.TxOptions) pgx.TxOptions

func WithIsoLevel(level pgx.TxIsoLevel) TxOptionFunc {
	return func(opt pgx.TxOptions) pgx.TxOptions {
		opt.IsoLevel = level
		return opt
	}
}

func (s *Storage) Tx(ctx context.Context, opts ...TxOptionFunc) (adaptors.ConnectionWrapper, error) {
	options := pgx.TxOptions{}
	for _, opt := range opts {
		options = opt(options)
	}

	tx, err := s.db.BeginTx(ctx, options)
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context, command func(ctx context.Context, q pgxtype.Querier) error) error {
		if err := command(ctx, tx); err != nil {
			_ = tx.Rollback(ctx)
			return err
		}

		return tx.Commit(ctx)
	}, nil
}
