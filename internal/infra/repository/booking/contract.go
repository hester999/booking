package booking

import (
	"booking/internal/domain"
	"context"
)

type bookingRepository interface {
	CreateBooking(ctx context.Context, booking *domain.Booking) error
	GetBooking(ctx context.Context, id string) (*domain.Booking, error)
	UpdateStatus(ctx context.Context, booking *domain.Booking) (*domain.Booking, error)
	GetExpiredBookings(ctx context.Context) ([]*domain.Booking, error)
	HasOverlappingBooking(ctx context.Context, booking *domain.Booking) (bool, error)
	GetBookingsByListingID(ctx context.Context, listingID string) ([]*domain.Booking, error)
}

type bookingCache interface {
	GetBooking(ctx context.Context, id string) (*domain.Booking, error)
	SetBooking(ctx context.Context, booking *domain.Booking) error
	DeleteBooking(ctx context.Context, id string) error
}
