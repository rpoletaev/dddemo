package payments

import (
	"context"

	"github.com/pkg/errors"
)

type PaymentMethod struct {
	Name string
	ID   uint64
}

var (
	ErrNotFound   = errors.New("not found")
	ErrBadRequest = errors.New("bad request")
)

type PaymentsBindingService interface {
	GetActivePaymentMethodForUser(ctx context.Context, userID uint64) (PaymentMethod, error)
	GetBindingForm(ctx context.Context, userID uint64) (string, error)
	Bind(ctx context.Context, pm PaymentMethod, userID uint64, orderID string) (string, error)
}

type PaymentsRepository interface {
	Save(p Payment) error
	GetByUUID(string) (Payment, error)
}
type BindCommand struct {
	storage PaymentsRepository
	service PaymentsBindingService
}

func (cmd BindCommand) Run(ctx context.Context, userID uint64) error {

}
