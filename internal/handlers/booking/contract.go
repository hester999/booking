package booking

import (
	command "booking/internal/application/command/booking"
	query "booking/internal/application/query/booking"
	"context"
)

type booking interface {
	CreateBooking(ctx context.Context, createBookingCommand command.CreateBookingCommand) (*query.BookingResult, error)
	Confirm(ctx context.Context, confirmBookingCommand command.ConfirmBookingCommand) error
	Cancel(ctx context.Context, cancelBookingCommand command.CancelBookingCommand) error
	//Expire(ctx context.Context, id string) error
	GetBooking(ctx context.Context, bookingID string) (*query.BookingResult, error)

	//GetExpiredBooking(ctx context.Context) (query.BookingsResult, error)
}
