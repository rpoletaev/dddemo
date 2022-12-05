package commands

import (
	"context"
	"subscription/internal/adaptors"
	"subscription/internal/adaptors/postgres"
	subscriptionsRepository "subscription/internal/adaptors/postgres/subscriptions"
	"subscription/internal/domain/subscriptions"

	"github.com/jackc/pgtype/pgxtype"
)

type UnitOfWorkFactory struct {
	storage *postgres.Storage
}

type UnitOfWork interface {
	Do(context.Context) error
}

func (f *UnitOfWorkFactory) SubscribeCmd(ctx context.Context, userID uint64) (UnitOfWork, error) {
	wrapper, err := f.storage.Tx(ctx)
	if err != nil {
		return nil, err
	}

	cmd := SubscribeCMD{
		connWrapper: wrapper,
		f: func(cxt context.Context, db pgxtype.Querier) error {
			subscriptionsRepo := subscriptionsRepository.New(db)
			subscribeCmd := subscriptions.NewSubscribeCmd(subscriptionsRepo, userID)
			return subscribeCmd.Run(ctx)
		},
	}

	return cmd, nil
}

type SubscribeCMD struct {
	f           func(ctx context.Context, q pgxtype.Querier) error
	connWrapper adaptors.ConnectionWrapper
}

func (cmd SubscribeCMD) Do(ctx context.Context) error {
	return cmd.connWrapper(ctx, cmd.f)
}
