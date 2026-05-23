package booking

import (
	"booking/internal/domain"
	"encoding/json"
	"time"
)

type BookingCreatedPayload struct {
	BookingID string `json:"booking_id"`
	UserID    string `json:"user_id"`
	ListingID string `json:"listing_id"`
	DateFrom  string `json:"date_from"`
	DateTo    string `json:"date_to"`
	Status    string `json:"status"`
}

func buildBookingPayload(snap domain.BookingSnapshot) ([]byte, error) {
	payloadDTO := BookingCreatedPayload{
		BookingID: snap.ID,
		UserID:    snap.UserID,
		ListingID: snap.ListingID,
		DateFrom:  snap.DateFrom.Format(time.RFC3339),
		DateTo:    snap.DateTo.Format(time.RFC3339),
		Status:    string(snap.Status),
	}

	return json.Marshal(payloadDTO)
}
