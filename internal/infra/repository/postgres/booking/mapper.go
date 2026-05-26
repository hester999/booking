package booking

import (
	"booking/apperr"
	"booking/internal/domain"
)

func domainBookingToDB(booking *domain.Booking) *BookingRepository {
	snap := booking.Snapshot()

	return &BookingRepository{
		ID:         snap.ID,
		UserID:     snap.UserID,
		ListingID:  snap.ListingID,
		DateFrom:   snap.DateFrom,
		DateTo:     snap.DateTo,
		Status:     string(snap.Status),
		ExpiresAt:  snap.ExpiredAt,
		CreatedAt:  snap.CreatedAt,
		UpdatedAt:  snap.UpdatedAt,
		TotalPrice: snap.PriceDecimal(),
	}
}

func repoModelToDomain(result *BookingRepository) (*domain.Booking, error) {
	convertedStatus, err := domain.ParseStatus(result.Status)
	if err != nil {
		return nil, apperr.InvalidStatusErr
	}

	price, err := domain.NewPriceFromDecimal(result.TotalPrice)
	if err != nil {
		return nil, err
	}

	snapBooking := domain.BookingSnapshot{
		ID:         result.ID,
		UserID:     result.UserID,
		ListingID:  result.ListingID,
		DateFrom:   result.DateFrom,
		DateTo:     result.DateTo,
		Status:     convertedStatus,
		ExpiredAt:  result.ExpiresAt,
		CreatedAt:  result.CreatedAt,
		UpdatedAt:  result.UpdatedAt,
		TotalPrice: price,
	}

	restoreBooking, err := snapBooking.RestoreDomain()
	if err != nil {
		return nil, err
	}
	return restoreBooking, nil
}
