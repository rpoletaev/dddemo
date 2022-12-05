package subscriptions

import (
	"context"
	"fmt"
	"strings"
	"subscription/internal/domain"
	"subscription/internal/domain/subscriptions"
	"time"

	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

type Repository struct {
	db pgxtype.Querier
}

func New(q pgxtype.Querier) Repository {
	return Repository{db: q}
}

const getSubscriptionQueryTmpl = "SELECT id, user_id, status, status_changed_at, is_need_to_prolong, active_since, active_until FROM subscriptions WHERE %s = $1"
const insertSubscriptionQuery = `INSERT INTO subscriptions (id, user_id, status, status_changed_at, is_need_to_prolong, active_since, active_until)
 VALUES (?, ?, ?, ?, ?, ?, ?)
 ON CONFLICT(id) DO UPDATE
 user_id=EXCLUDED.user_id, 
 status=EXCLUDED.status, 
 status_changed_at=EXCLUDED.status_changed_at, 
 is_need_to_prolong=EXCLUDED.is_need_to_prolong, 
 active_since=EXCLUDED.active_since, 
 active_until=EXCLUDED.active_until
 `
const forUpdateTmpl = "FOR UPDATE SKIP LOCKED"

func (r Repository) NewID(ctx context.Context) (uint64, error) {
	var id uint64
	err := r.db.QueryRow(ctx, "SELECT nextval('subscriptions_id_seq')").Scan(&id)
	return id, err
}

func (r Repository) GetByID(ctx context.Context, id uint64) (subscriptions.Subscription, error) {
	return r.get(ctx, "id", id, false)
}

func (r Repository) GetByIDForUpdate(ctx context.Context, id uint64) (subscriptions.Subscription, error) {
	return r.get(ctx, "id", id, true)
}

func (r Repository) GetByUserID(ctx context.Context, id uint64) (subscriptions.Subscription, error) {
	return r.get(ctx, "user_id", id, false)
}

func (r Repository) Save(ctx context.Context, sub subscriptions.Subscription) error {
	_, err := r.db.Exec(
		ctx,
		insertSubscriptionQuery,
		sub.Id(),
		sub.UserID(),
		sub.Status(),
		sub.StatusChangedAt(),
		sub.IsNeedToProlong(),
		sub.ActiveSince(),
		sub.ActiveUntil(),
	)

	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			return domain.ErrAlreadyExists
		}
		return errors.Wrap(domain.ErrDBError, err.Error())
	}

	return nil
}

func (r Repository) get(ctx context.Context, field string, value uint64, forUpdate bool) (subscriptions.Subscription, error) {
	query := fmt.Sprintf(getSubscriptionQueryTmpl, field)
	if forUpdate {
		query = strings.Join([]string{query, forUpdateTmpl}, " ")
	}

	row := r.db.QueryRow(ctx, query, value)
	var id uint64
	var userID uint64
	var status string
	var statusChangedAt time.Time
	var isNeedToProlong bool
	var activeSince *time.Time
	var activeUntil *time.Time

	if err := row.Scan(&id, &userID, &status, &statusChangedAt, &isNeedToProlong, &activeSince, &activeUntil); err != nil {
		if err == pgx.ErrNoRows {
			return subscriptions.Subscription{}, domain.ErrNotFound
		}

		return subscriptions.Subscription{}, errors.Wrap(domain.ErrDBError, err.Error())
	}

	return subscriptions.New(
		id,
		userID,
		status,
		statusChangedAt,
		isNeedToProlong,
		activeSince,
		activeUntil)
}
