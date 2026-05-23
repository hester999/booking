package booking

import "time"

type BookingResult struct {
	ID         string
	ListingID  string
	Status     string
	DateTo     time.Time
	DateFrom   time.Time
	ExpiredAt  *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	TotalPrice string
}

type BookingsResult []*BookingResult
