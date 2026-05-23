package cache

import (
	"booking/apperr"
	"booking/internal/domain"
	"time"
)

type CacheBooking struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	ListingID  string     `json:"listing_id"`
	DateFrom   time.Time  `json:"date_from"`
	DateTo     time.Time  `json:"date_to"`
	Status     string     `json:"status"`
	ExpiresAt  *time.Time `json:"expires_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	TotalPrice string     `json:"total_price"`
}

func (c *CacheBooking) DomainToDTOCreate(booking *domain.Booking) {
	snap := booking.Snapshot()

	c.ID = snap.ID
	c.UserID = snap.UserID
	c.ListingID = snap.ListingID
	c.DateFrom = snap.DateFrom
	c.DateTo = snap.DateTo
	c.Status = string(snap.Status)
	c.ExpiresAt = snap.ExpiredAt
	c.CreatedAt = snap.CreatedAt
	c.UpdatedAt = snap.UpdatedAt
	c.TotalPrice = snap.PriceToString()

}

func (c *CacheBooking) CacheToDomain() (*domain.Booking, error) {
	convertedStatus, err := domain.ParseStatus(c.Status)
	if err != nil {
		return nil, apperr.InvalidStatusErr
	}
	snapBooking := domain.BookingSnapshot{
		ID:        c.ID,
		UserID:    c.UserID,
		ListingID: c.ListingID,
		DateFrom:  c.DateFrom,
		DateTo:    c.DateTo,
		Status:    convertedStatus,
		ExpiredAt: c.ExpiresAt,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
	err = snapBooking.StringToDecimal(c.TotalPrice)
	if err != nil {
		return nil, err
	}
	restoreBooking, err := snapBooking.RestoreDomain()
	if err != nil {
		return nil, err
	}
	return restoreBooking, nil
}
