package dto

import (
	command "booking/internal/application/command/booking"
	"strings"
	"time"
)

type CreateBookingRequest struct {
	UserID    string `json:"user_id"`
	ListingID string `json:"listing_id"`
	DateFrom  string `json:"date_from"`
	DateTo    string `json:"date_to"`
}

func (req *CreateBookingRequest) ToCreateCommand() (command.CreateBookingCommand, error) {
	dateFrom, err := parseDate(req.DateFrom)
	if err != nil {
		return command.CreateBookingCommand{}, err
	}
	dateTo, err := parseDate(req.DateTo)
	if err != nil {
		return command.CreateBookingCommand{}, err
	}

	return command.CreateBookingCommand{
		UserID:    req.UserID,
		ListingID: req.ListingID,
		DateFrom:  dateFrom,
		DateTo:    dateTo,
	}, nil
}

func parseDate(date string) (time.Time, error) {
	return time.Parse("2006-01-02", strings.TrimSpace(date))
}

type BookingConfirmRequest struct {
	ID string `json:"id"`
}

type BookingCancelRequest struct {
	ID string `json:"id"`
}
