package booking

import (
	"booking/apperr"
	"booking/internal/domain"
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type Booking struct {
	db      *sqlx.DB
	tx      *trmsqlx.CtxGetter
	builder sq.StatementBuilderType
}

const (
	expired = "expired"
)

func NewBookingRepo(db *sqlx.DB, tx *trmsqlx.CtxGetter) *Booking {
	return &Booking{
		db:      db,
		tx:      tx,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (b *Booking) CreateBooking(ctx context.Context, booking *domain.Booking) error {
	tx := b.tx.DefaultTrOrDB(ctx, b.db)

	dto := domainBookingToDB(booking)

	builder := b.builder.Insert("booking").
		Columns("id", "user_id", "listing_id", "status", "date_from", "date_to", "expires_at", "total_price").
		Values(dto.ID, dto.UserID, dto.ListingID, dto.Status, dto.DateFrom, dto.DateTo, dto.ExpiresAt, dto.TotalPrice).
		Suffix("returning id,user_id,listing_id,date_from,date_to,status,expires_at,created_at,updated_at,total_price")

	sqlQuery, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("could not generate sql query: %w", err)
	}

	_, err = tx.ExecContext(ctx, sqlQuery, args...)
	if err != nil {
		return fmt.Errorf("could not insert booking: %w", err)
	}
	return nil
}

func (b *Booking) UpdateStatus(ctx context.Context, booking *domain.Booking) (*domain.Booking, error) {
	tx := b.tx.DefaultTrOrDB(ctx, b.db)

	dto := domainBookingToDB(booking)

	query := `update booking set status = $1, expires_at = $2, updated_at = now() where id = $3 returning id, user_id, listing_id, date_from, date_to, status, expires_at, created_at, updated_at`

	err := tx.GetContext(ctx, dto, query, dto.Status, dto.ExpiresAt, dto.ID)

	if err != nil {
		return nil, fmt.Errorf("could not update status: %w", err)
	}

	domainModel, err := repoModelToDomain(dto)
	if err != nil {
		return nil, err
	}
	return domainModel, nil
}

func (b *Booking) GetExpiredBookings(ctx context.Context) ([]*domain.Booking, error) {
	tx := b.tx.DefaultTrOrDB(ctx, b.db)

	query := `select id, user_id, listing_id, date_from, date_to, status, expires_at, created_at, updated_at, total_price from booking where expires_at < now() and status != $1` //todo тут нужен статус или нет?

	var result []BookingRepository

	err := tx.SelectContext(ctx, &result, query, expired)

	if err != nil {
		return nil, fmt.Errorf("could not get expired bookings: %w", err)
	}

	var bookings []*domain.Booking

	for _, v := range result {
		tmp, err := repoModelToDomain(&v)
		if err != nil {
			return nil, err
		}

		bookings = append(bookings, tmp)
	}
	return bookings, err
}

func (b *Booking) GetBooking(ctx context.Context, id string) (*domain.Booking, error) {
	tx := b.tx.DefaultTrOrDB(ctx, b.db)
	query := `select id,user_id,listing_id,date_from,date_to,status,expires_at,created_at,updated_at,total_price from booking where id = $1`

	var result BookingRepository
	err := tx.GetContext(ctx, &result, query, id)

	if err != nil {
		if _, ok := errors.AsType[*pgconn.PgError](err); ok {
			return nil, apperr.NotFoundErr
		}
		return nil, fmt.Errorf("could not get booking: %w", err)
	}

	domainModel, err := repoModelToDomain(&result)
	if err != nil {
		return nil, err
	}
	return domainModel, nil
}

func (b *Booking) GetBookingsByListingID(ctx context.Context, listingID string) ([]*domain.Booking, error) {
	var result []BookingRepository
	tx := b.tx.DefaultTrOrDB(ctx, b.db)

	query := `select id,user_id,listing_id,date_from,date_to,status,expires_at,created_at,updated_at,total_price from booking where listing_id = $1 and status = 'confirmed'`

	err := tx.GetContext(ctx, &result, query, listingID)
	if err != nil {
		if _, ok := errors.AsType[*pgconn.PgError](err); ok {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, apperr.NotFoundErr
			}
		}
	}

	var bookings []*domain.Booking

	for _, v := range result {
		tmp, err := repoModelToDomain(&v)
		if err != nil {
			return nil, err
		}

		bookings = append(bookings, tmp)
	}
	return bookings, err
}

func (b *Booking) HasOverlappingBooking(ctx context.Context, booking *domain.Booking) (bool, error) {
	tx := b.tx.DefaultTrOrDB(ctx, b.db)
	dto := domainBookingToDB(booking)

	query := `
		select exists (
			select 1
			from booking
			where listing_id = $1
			  and status in ('pending', 'confirmed')
			  and date_from < $3
			  and $2 < date_to
		);
	`

	var exists bool
	if err := tx.GetContext(ctx, &exists, query, dto.ListingID, dto.DateFrom, dto.DateTo); err != nil {
		return false, err
	}

	return exists, nil
}
