package domain

import (
	"booking/apperr"
	"time"
)

type Day int64

type Period struct {
	dateFrom time.Time
	dateTo   time.Time
}

func NewPeriod(dateFrom, dateTo time.Time) (*Period, error) {
	from := normalizeDate(dateFrom)
	to := normalizeDate(dateTo)

	if !to.After(from) {
		return nil, apperr.InvalidPeriodErr
	}

	return &Period{
		dateFrom: from,
		dateTo:   to,
	}, nil
}
func normalizeDate(t time.Time) time.Time {
	u := t.UTC()
	return time.Date(u.Year(), u.Month(), u.Day(), 0, 0, 0, 0, time.UTC)
}
