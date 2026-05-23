package expired

import (
	query "booking/internal/application/query/booking"
	"context"
)

type bookingUsecases interface {
	GetExpiredBooking(ctx context.Context) (query.BookingsResult, error)
	Expire(ctx context.Context, id string) error
}

type tx interface {
	Run(ctx context.Context, fn func(ctx context.Context) error) error
}
