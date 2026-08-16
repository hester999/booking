package booking

import (
	"booking/internal/domain"
	"booking/pkg"
	"context"
)

type BookingRepo struct {
	pg    bookingRepository
	cache bookingCache
	log   pkg.Logger
}

func NewCachedBookingRepo(pg bookingRepository, cache bookingCache, log pkg.Logger) *BookingRepo {
	return &BookingRepo{
		pg:    pg,
		cache: cache,
		log:   log,
	}
}

func (b *BookingRepo) CreateBooking(ctx context.Context, newBooking *domain.Booking) error {
	return b.pg.CreateBooking(ctx, newBooking)
}

func (b *BookingRepo) GetBooking(ctx context.Context, id string) (*domain.Booking, error) {
	result, err := b.cache.GetBooking(ctx, id)
	if err == nil && result != nil {
		return result, nil
	}

	if err != nil {
		b.log.Error("cache get booking failed", "error", err)
		err = nil
	}

	booking, err := b.pg.GetBooking(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := b.cache.SetBooking(ctx, booking); err != nil {
		b.log.Error("cache set booking failed", "error", err)
	}

	return booking, nil
}

func (b *BookingRepo) UpdateStatus(ctx context.Context, booking *domain.Booking) (*domain.Booking, error) {
	result, err := b.pg.UpdateStatus(ctx, booking)

	if err != nil {
		return nil, err
	}

	if err := b.cache.DeleteBooking(ctx, result.ID()); err != nil {
		b.log.Error("cache delete booking failed", "error", err)
	}

	return result, nil
}

func (b *BookingRepo) GetExpiredBookings(ctx context.Context) ([]*domain.Booking, error) {
	return b.pg.GetExpiredBookings(ctx)
}

func (b *BookingRepo) GetBookingsByListingID(ctx context.Context, listingID string) ([]*domain.Booking, error) {
	return b.pg.GetBookingsByListingID(ctx, listingID)
}

func (b *BookingRepo) HasOverlappingBooking(ctx context.Context, booking *domain.Booking) (bool, error) {
	return b.pg.HasOverlappingBooking(ctx, booking)
}

func (b *BookingRepo) LockListing(ctx context.Context, listingID string) error {
	return b.pg.LockListing(ctx, listingID)
}
