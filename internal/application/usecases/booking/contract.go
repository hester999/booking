package booking

import (
	"booking/internal/domain"
	"context"
)

type tx interface {
	Run(ctx context.Context, fn func(ctx context.Context) error) error
}

type eventRepo interface {
	CreateEvent(ctx context.Context, newEvent *domain.Event) error
}

type bookingRepo interface {
	CreateBooking(ctx context.Context, newBooking *domain.Booking) error
	GetBooking(ctx context.Context, id string) (*domain.Booking, error)
	UpdateStatus(ctx context.Context, booking *domain.Booking) (*domain.Booking, error)
	GetExpiredBookings(ctx context.Context) ([]*domain.Booking, error)
	GetBookingsByListingID(ctx context.Context, listingID string) ([]*domain.Booking, error)
	HasOverlappingBooking(ctx context.Context, booking *domain.Booking) (bool, error)
}
