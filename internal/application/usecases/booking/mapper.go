package booking

import (
	query "booking/internal/application/query/booking"
	"booking/internal/domain"
)

func DomainToResult(booking domain.BookingSnapshot) *query.BookingResult {
	return &query.BookingResult{
		ID:         booking.ID,
		ListingID:  booking.ListingID,
		Status:     string(booking.Status),
		DateTo:     booking.DateTo,
		DateFrom:   booking.DateFrom,
		ExpiredAt:  booking.ExpiredAt,
		CreatedAt:  booking.CreatedAt,
		UpdatedAt:  booking.UpdatedAt,
		TotalPrice: booking.PriceToString(),
	}
}

func DomainsToResult(bookings []domain.BookingSnapshot) query.BookingsResult {
	result := make(query.BookingsResult, 0, len(bookings))
	for _, booking := range bookings {
		result = append(result, DomainToResult(booking))
	}
	return result
}
