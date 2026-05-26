package booking

import (
	"time"

	"github.com/shopspring/decimal"
)

type BookingRepository struct {
	ID         string          `db:"id"`
	UserID     string          `db:"user_id"`
	ListingID  string          `db:"listing_id"`
	DateFrom   time.Time       `db:"date_from"`
	DateTo     time.Time       `db:"date_to"`
	Status     string          `db:"status"`
	ExpiresAt  *time.Time      `db:"expires_at"`
	CreatedAt  time.Time       `db:"created_at"`
	UpdatedAt  time.Time       `db:"updated_at"`
	TotalPrice decimal.Decimal `db:"total_price"`
}
