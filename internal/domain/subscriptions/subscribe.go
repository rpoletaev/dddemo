package subscriptions

import (
	"context"
	"time"
)

type SubscriptionRepository interface {
	NewID(ctx context.Context) (uint64, error)
	GetByID(ctx context.Context, id uint64) (Subscription, error)
	GetByIDForUpdate(ctx context.Context, id uint64) (Subscription, error)
	GetByUserID(ctx context.Context, id uint64) (Subscription, error)
	Save(ctx context.Context, sub Subscription) error
}

type Subscribe struct {
	subRepo SubscriptionRepository
	params  subscribeParams
}

type subscribeParams struct {
	userID uint64
}

func NewSubscribeCmd(repo SubscriptionRepository, userID uint64) Subscribe {
	return Subscribe{
		subRepo: repo,
		params:  subscribeParams{userID: userID},
	}
}

func (s Subscribe) Run(ctx context.Context) error {
	id, err := s.subRepo.NewID(ctx)
	if err != nil {
		return err
	}

	sub, err := New(
		id,
		s.params.userID,
		StatusWaiting,
		time.Now(),
		false,
		nil,
		nil,
	)

	if err != nil {
		return err
	}

	return s.subRepo.Save(ctx, sub)
}
