package booking

import "time"

type CreateBookingCommand struct {
	UserID    string
	ListingID string
	DateFrom  time.Time
	DateTo    time.Time
}

type ConfirmBookingCommand struct {
	ID string
}

type CancelBookingCommand struct {
	ID string
}
