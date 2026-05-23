package domain

import (
	"booking/apperr"
	"booking/utils"
	"time"
)

type Status string

type Booking struct {
	id         string
	userID     string
	listingID  string
	status     Status
	period     Period
	totalPrice Price
	expiredAt  *time.Time
	createdAt  time.Time
	updatedAt  time.Time
}

const (
	StatusPending   Status = "pending"
	StatusConfirmed Status = "confirmed"
	StatusCancelled Status = "cancelled"
	StatusExpired   Status = "expired"
)

func NewBooking(userID string, listingID string, dateFrom, dateTo, expiredAt time.Time) (*Booking, error) {
	var err error

	if userID == "" {
		return nil, apperr.InvalidUserIDErr
	}

	if listingID == "" {
		return nil, apperr.InvalidListingID
	}

	period, err := NewPeriod(dateFrom, dateTo)
	if err != nil {
		return nil, err
	}
	today := normalizeDate(time.Now())

	if period.dateFrom.Before(today) {
		return nil, apperr.InvalidDateErr
	}

	newBooking := &Booking{
		id:        utils.GetUUID(),
		userID:    userID,
		listingID: listingID,
		expiredAt: &expiredAt,
		status:    StatusPending,
		period:    *period,
	}

	return newBooking, nil

}

func (b *Booking) ID() string {
	return b.id
}

func (b *Booking) Snapshot() BookingSnapshot {
	return BookingSnapshot{
		ID:         b.id,
		UserID:     b.userID,
		ListingID:  b.listingID,
		Status:     b.status,
		DateFrom:   b.period.dateFrom,
		DateTo:     b.period.dateTo,
		ExpiredAt:  b.expiredAt,
		CreatedAt:  b.createdAt,
		UpdatedAt:  b.updatedAt,
		TotalPrice: b.totalPrice,
	}
}

func (b *Booking) Cancel() error {
	switch b.status {
	case StatusExpired:
		return apperr.StatusExpire
	case StatusConfirmed:
		return apperr.StatusConfirmed
	case StatusCancelled:
		return nil
	case StatusPending:
		b.status = StatusCancelled
		b.expiredAt = nil
		return nil
	default:
		return apperr.UnsupportedStatus
	}

}

func (b *Booking) Confirm() error {

	switch b.status {
	case StatusCancelled:
		return apperr.StatusCancel
	case StatusExpired:
		return apperr.StatusExpire
	case StatusConfirmed:
		return nil
	case StatusPending:
		b.status = StatusConfirmed
		b.expiredAt = nil
		return nil
	default:
		return apperr.UnsupportedStatus
	}

}

func (b *Booking) Expired() error {

	switch b.status {
	case StatusConfirmed:
		return apperr.StatusConfirmed
	case StatusCancelled:
		return apperr.StatusCancel
	case StatusExpired:
		return nil
	case StatusPending:
		if b.expiredAt == nil {
			return apperr.InvalidDateErr
		}
		b.status = StatusExpired
		return nil
	default:
		return apperr.UnsupportedStatus
	}
}

func (b *Booking) IsExpired(now time.Time) bool {
	if b.expiredAt == nil {
		return false
	}
	if now.After(*b.expiredAt) {
		return true
	}
	return false
}

func (b *Booking) UpdatePeriod(dateFrom, dateTo time.Time) error {
	if b.status == StatusCancelled || b.status == StatusExpired {
		return apperr.InvalidPeriodErr
	}
	period, err := NewPeriod(dateFrom, dateTo)
	if err != nil {
		return err
	}
	b.period = *period
	return nil
}

func (b *Booking) CheckIntersection(startDate, endDate time.Time) bool {
	period, err := NewPeriod(startDate, endDate)
	if err != nil {
		return false
	}

	existFrom := b.period.dateFrom
	existTo := b.period.dateTo

	newFrom := period.dateFrom
	newTo := period.dateTo

	oldIsAfterNew := !existFrom.Before(newTo)
	oldIsBeforeNew := !existTo.After(newFrom)

	return !(oldIsAfterNew || oldIsBeforeNew)
}

func (b *Booking) CalculatePrice(price string) error {

	total, err := NewTotalPrice(price)
	if err != nil {
		return err
	}
	b.totalPrice = total
	return nil
}

func ParseStatus(value string) (Status, error) {
	switch value {
	case string(StatusPending):
		return StatusPending, nil
	case string(StatusConfirmed):
		return StatusConfirmed, nil
	case string(StatusCancelled):
		return StatusCancelled, nil
	case string(StatusExpired):
		return StatusExpired, nil
	default:
		return "", apperr.InvalidStatusErr
	}
}
