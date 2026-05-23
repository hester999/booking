package domain

import (
	"time"
)

type BookingSnapshot struct {
	ID         string
	UserID     string
	ListingID  string
	Status     Status
	DateFrom   time.Time
	DateTo     time.Time
	ExpiredAt  *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	TotalPrice Price
}

func (s *BookingSnapshot) RestoreDomain() (*Booking, error) {
	period, err := NewPeriod(s.DateFrom, s.DateTo)
	if err != nil {
		return nil, err
	}
	return &Booking{
		id:         s.ID,
		userID:     s.UserID,
		listingID:  s.ListingID,
		status:     s.Status,
		expiredAt:  s.ExpiredAt,
		createdAt:  s.CreatedAt,
		updatedAt:  s.UpdatedAt,
		period:     *period,
		totalPrice: s.TotalPrice,
	}, nil
}

func (s *BookingSnapshot) PriceToString() string {
	return s.TotalPrice.totalPrice.String()
}

func (s *BookingSnapshot) StringToDecimal(price string) error {
	newPrice, err := NewTotalPrice(price)
	if err != nil {
		return err
	}
	s.TotalPrice = newPrice

	return nil
}
