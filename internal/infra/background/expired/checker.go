package expired

import (
	"booking/pkg"
	"context"
	"fmt"
	"time"
)

type ExpiredChecker struct {
	booking bookingUsecases
	tx      tx
	delay   int64
	errChan chan error
	logger  pkg.Logger
}

func NewExpiredChecker(booking bookingUsecases, delay int64, tx tx, logger pkg.Logger) *ExpiredChecker {
	return &ExpiredChecker{
		booking: booking,
		tx:      tx,
		delay:   delay,
		errChan: make(chan error, 10),
		logger:  logger,
	}
}

func (e *ExpiredChecker) GetErrChan() <-chan error {
	return e.errChan
}

func (e *ExpiredChecker) Start(ctx context.Context) error {
	ticker := time.NewTicker(time.Duration(e.delay) * time.Second)
	defer close(e.errChan)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			expiredBookings, err := e.booking.GetExpiredBooking(ctx)
			if err != nil {
				e.errChan <- err
				continue
			}
			for _, booking := range expiredBookings {
				err = e.booking.Expire(ctx, booking.ID)
				if err != nil {
					e.errChan <- err
					continue
				}
				e.logger.Info(fmt.Sprintf("Expired Booking: %s", booking.ID))
			}

		}
	}
}
