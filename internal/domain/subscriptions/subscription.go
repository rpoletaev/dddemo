package subscriptions

import (
	"time"

	"github.com/pkg/errors"
)

func isStatusValid(status string) bool {

	if status == StatusActive || status == StatusInactive || status == StatusWaiting {
		return true
	}

	return false
}

const (
	StatusActive   = "active"
	StatusInactive = "inactive"
	StatusWaiting  = "waiting"
)

var (
	ErrAlreadyInStatus = errors.New("already in status")
	ErrUnknownStatus   = errors.New("unknown status")
)

type Subscription struct {
	id              uint64
	userID          uint64
	status          string
	statusChangedAt time.Time
	isNeedToProlong bool
	activeSince     *time.Time
	activeUntil     *time.Time
}

func (s Subscription) Id() uint64 {
	return s.id
}
func (s Subscription) UserID() uint64 {
	return s.userID
}
func (s Subscription) Status() string {
	return s.status
}
func (s Subscription) StatusChangedAt() time.Time {
	return s.statusChangedAt
}
func (s Subscription) IsNeedToProlong() bool {
	return s.isNeedToProlong
}
func (s Subscription) ActiveSince() *time.Time {
	return s.activeSince
}
func (s Subscription) ActiveUntil() *time.Time {
	return s.activeUntil
}

// New returns instantiated subscription object
// if wrong status passed in then returns ErrUnknownStatus
func New(
	id uint64,
	userID uint64,
	status string,
	statusChangedAt time.Time,
	isNeedToProlong bool,
	activeSince *time.Time,
	activeUntil *time.Time,
) (Subscription, error) {
	if !isStatusValid(status) {
		return Subscription{}, errors.Wrap(ErrUnknownStatus, "status")
	}

	s := Subscription{
		id:              id,
		userID:          userID,
		status:          status,
		statusChangedAt: statusChangedAt,
		isNeedToProlong: isNeedToProlong,
		activeSince:     activeSince,
		activeUntil:     activeUntil,
	}

	return s, nil
}

// Activate makes subscription to active state.
// In this state all subscription's advantages is available for user.
// Autoprolongation becomes enabled.&
func (s *Subscription) Activate(since, until time.Time) error {
	if s.status == StatusActive {
		return ErrAlreadyInStatus
	}

	s.status = StatusActive
	s.isNeedToProlong = true
	s.statusChangedAt = time.Now()
	s.activeSince = &since
	s.activeUntil = &until
	return nil
}

// Deactivate turns off the subscription.
// Turns off next prolongation.
// All advantages becomes unavailable
func (s *Subscription) Deactivate() error {
	if s.status == StatusInactive {
		return ErrAlreadyInStatus
	}

	s.status = StatusInactive
	s.isNeedToProlong = false
	s.statusChangedAt = time.Now()
	s.activeSince = nil
	s.activeUntil = nil
	return nil
}

// DisableProlongation turns off next prolongation.
// All advantagest are still available.
func (s *Subscription) DisableProlongation() {
	s.isNeedToProlong = false
}
