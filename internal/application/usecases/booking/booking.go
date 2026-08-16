package booking

import (
	"booking/apperr"
	command "booking/internal/application/command/booking"
	query "booking/internal/application/query/booking"
	"booking/internal/domain"
	"context"
	"fmt"
	"time"
)

type Booking struct {
	tx             tx
	bookingRepo    bookingRepo
	eventRepo      eventRepo
	listingGateway listingGateway
}

func NewBookingUsecase(bookingRepo bookingRepo, eventRepo eventRepo, listingGateway listingGateway, tx tx) *Booking {
	return &Booking{
		bookingRepo:    bookingRepo,
		eventRepo:      eventRepo,
		tx:             tx,
		listingGateway: listingGateway,
	}
}

const (
	op = "booking"
)

// todo  для создания  брони необходимо проверять наличие пользователя и объявления в других сервисах, например цену за 1 сутки
//  Получать инфу, считать общий прайс в бизнес логике.
//  Проверять наличие пользователя, его права, персональную скидку
//  Из сервиса  listing приходит моделька, в ней есть минимальная длительность брони, учесть в проверке dateFrom/dateTo. Пример 3 дня, если  длительность брони 2 дня, то ошибка

func (b *Booking) CreateBooking(ctx context.Context, createBookingCommand command.CreateBookingCommand) (*query.BookingResult, error) {
	var bookingSnap domain.BookingSnapshot
	listing, gateErr := b.listingGateway.GetListing(ctx, createBookingCommand.ListingID)
	if gateErr != nil {
		return nil, fmt.Errorf("[%s] get listing from external service error: %w", op, gateErr)
	}

	err := b.tx.Run(ctx, func(ctx context.Context) error {

		domainBooking, err := domain.NewBooking(
			createBookingCommand.UserID,
			createBookingCommand.ListingID,
			createBookingCommand.DateFrom,
			createBookingCommand.DateTo,
			time.Now().UTC().Add(3*time.Hour), // todo че делать с этим
		)

		if err != nil {
			return fmt.Errorf("[%s] create booking error  %w", op, err)
		}

		hasOverlap, err := b.bookingRepo.HasOverlappingBooking(ctx, domainBooking)
		if err != nil {
			return err
		}

		if hasOverlap {
			return apperr.BookingAlreadyExistsErr
		}

		err = domainBooking.ValidateMinimalBookingDuration(listing.MinimumBookingDuration)
		if err != nil {
			return err
		}

		err = domainBooking.CalculatePrice(listing.Price, listing.MinimumBookingDuration)
		if err != nil {
			return err
		}

		err = b.bookingRepo.LockListing(ctx, domainBooking.ListingID())

		err = b.bookingRepo.CreateBooking(ctx, domainBooking)

		if err != nil {
			return fmt.Errorf("[%s] create booking error  %w", op, err)
		}

		bookingSnap = domainBooking.Snapshot()

		payload, err := buildBookingPayload(bookingSnap)
		if err != nil {
			return fmt.Errorf("[%s] marshal event error  %w", op, err)
		}

		newEvent := domain.NewBookingCreatedEvent(bookingSnap.ID, payload)

		err = b.eventRepo.CreateEvent(ctx, newEvent)
		if err != nil {
			return fmt.Errorf("[%s] create event error  %w", op, err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return DomainToResult(bookingSnap), nil
}

func (b *Booking) Confirm(ctx context.Context, confirmBookingCommand command.ConfirmBookingCommand) error {

	err := b.tx.Run(ctx, func(ctx context.Context) error {

		booking, err := b.bookingRepo.GetBooking(ctx, confirmBookingCommand.ID)
		if err != nil {
			return fmt.Errorf("[%s] confirm booking error  %w", op, err)

		}

		if booking.IsExpired(time.Now().UTC()) {

			return fmt.Errorf("[%s] confirm booking error  %w", op, apperr.StatusExpire)

		}

		err = booking.Confirm()
		if err != nil {
			return fmt.Errorf("[%s] confirm booking error  %w", op, err)
		}

		_, err = b.bookingRepo.UpdateStatus(ctx, booking)
		if err != nil {
			return fmt.Errorf("[%s] confirm booking error  %w", op, err)
		}

		bookingSnap := booking.Snapshot()

		payload, err := buildBookingPayload(bookingSnap)
		if err != nil {
			return fmt.Errorf("[%s] marshal event error  %w", op, err)
		}

		newEvent := domain.NewBookingConfirmedEvent(bookingSnap.ID, payload)

		err = b.eventRepo.CreateEvent(ctx, newEvent)
		if err != nil {
			return fmt.Errorf("[%s] create event error  %w", op, err)
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (b *Booking) Cancel(ctx context.Context, cancelBookingCommand command.CancelBookingCommand) error {
	err := b.tx.Run(ctx, func(ctx context.Context) error {
		booking, err := b.bookingRepo.GetBooking(ctx, cancelBookingCommand.ID)
		if err != nil {
			return fmt.Errorf("[%s] cancel booking error  %w", op, err)
		}

		if booking.IsExpired(time.Now().UTC()) {
			return fmt.Errorf("[%s] confirm booking error  %w", op, apperr.StatusExpire)
		}

		err = booking.Cancel()
		if err != nil {
			return fmt.Errorf("[%s] cancel booking error  %w", op, err)

		}

		_, err = b.bookingRepo.UpdateStatus(ctx, booking)
		if err != nil {
			return fmt.Errorf("[%s] cancel booking error  %w", op, err)
		}

		bookingSnap := booking.Snapshot()

		payload, err := buildBookingPayload(bookingSnap)
		if err != nil {
			return fmt.Errorf("[%s] marshal event error  %w", op, err)
		}

		newEvent := domain.NewBookingCancelledEvent(bookingSnap.ID, payload)

		err = b.eventRepo.CreateEvent(ctx, newEvent)
		if err != nil {
			return fmt.Errorf("[%s] create event error  %w", op, err)
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *Booking) Expire(ctx context.Context, id string) error {
	err := b.tx.Run(ctx, func(ctx context.Context) error {

		booking, err := b.bookingRepo.GetBooking(ctx, id)
		if err != nil {
			return fmt.Errorf("[%s] expire booking error  %w", op, err)
		}
		if !booking.IsExpired(time.Now().UTC()) {
			return fmt.Errorf("[%s] expire booking error:  %w", op, apperr.NotExpiredErr)
		}

		err = booking.Expired()
		if err != nil {
			return fmt.Errorf("[%s] expire booking error  %w", op, err)
		}
		_, err = b.bookingRepo.UpdateStatus(ctx, booking)
		if err != nil {
			return fmt.Errorf("[%s] expire booking error  %w", op, err)
		}

		bookingSnap := booking.Snapshot()

		payload, err := buildBookingPayload(bookingSnap)
		if err != nil {
			return fmt.Errorf("[%s] marshal event error  %w", op, err)
		}

		newEvent := domain.NewBookingExpiredEvent(bookingSnap.ID, payload)

		err = b.eventRepo.CreateEvent(ctx, newEvent)
		if err != nil {
			return fmt.Errorf("[%s] create event error  %w", op, err)
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *Booking) GetBooking(ctx context.Context, bookingID string) (*query.BookingResult, error) {
	var bookingSnap domain.BookingSnapshot
	err := b.tx.Run(ctx, func(ctx context.Context) error {
		booking, err := b.bookingRepo.GetBooking(ctx, bookingID)
		if err != nil {
			return fmt.Errorf("[%s] get booking error  %w", op, err)
		}

		bookingSnap = booking.Snapshot()
		return nil
	})
	if err != nil {
		return nil, err
	}
	return DomainToResult(bookingSnap), nil
}

func (b *Booking) GetExpiredBooking(ctx context.Context) (query.BookingsResult, error) {
	var bookingSnap []domain.BookingSnapshot

	err := b.tx.Run(ctx, func(ctx context.Context) error {
		bookings, err := b.bookingRepo.GetExpiredBookings(ctx)
		if err != nil {
			return fmt.Errorf("[%s] get expired bookings error:  %w", op, err)
		}

		for _, booking := range bookings {
			bookingSnap = append(bookingSnap, booking.Snapshot())
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	res := DomainsToResult(bookingSnap)
	return res, nil
}
