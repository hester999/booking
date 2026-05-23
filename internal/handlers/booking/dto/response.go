package dto

import (
	"booking/internal/application/query/booking"
	"time"
)

type CreateBookingResponse struct {
	ID         string     `json:"id"`
	ListingID  string     `json:"listing_id"`
	Status     string     `json:"status"`
	DateTo     time.Time  `json:"date_to"`
	DateFrom   time.Time  `json:"date_from"`
	ExpiredAt  *time.Time `json:"expired_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	TotalPrice string     `json:"total_price"`
}

func (c *CreateBookingResponse) CreateQueryResultToResponse(result *booking.BookingResult) {
	c.ID = result.ID
	c.ListingID = result.ListingID
	c.Status = result.Status
	c.DateTo = result.DateTo
	c.DateFrom = result.DateFrom
	c.ExpiredAt = result.ExpiredAt
	c.CreatedAt = result.CreatedAt
	c.UpdatedAt = result.UpdatedAt
	c.TotalPrice = result.TotalPrice
}

type GetBookingResponse struct {
	ID         string     `json:"id"`
	ListingID  string     `json:"listing_id"`
	Status     string     `json:"status"`
	DateTo     time.Time  `json:"date_to"`
	DateFrom   time.Time  `json:"date_from"`
	ExpiredAt  *time.Time `json:"expired_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	TotalPrice string     `json:"total_price"`
}

func (c *GetBookingResponse) CreateQueryResultToResponse(result *booking.BookingResult) {
	c.ID = result.ID
	c.ListingID = result.ListingID
	c.Status = result.Status
	c.DateTo = result.DateTo
	c.DateFrom = result.DateFrom
	c.ExpiredAt = result.ExpiredAt
	c.CreatedAt = result.CreatedAt
	c.UpdatedAt = result.UpdatedAt
	c.TotalPrice = result.TotalPrice
}
