package adaptors

import (
	"context"

	"github.com/jackc/pgtype/pgxtype"
)

type ConnectionWrapper func(ctx context.Context, command func(ctx context.Context, q pgxtype.Querier) error) error
