package domain

import (
	"booking/apperr"

	"github.com/shopspring/decimal"
)

type Price struct {
	totalPrice decimal.Decimal
}

func NewTotalPrice(totalPrice string) (Price, error) {
	price, err := decimal.NewFromString(totalPrice)

	if err != nil {
		return Price{}, err
	}
	if price == decimal.Zero || price.IsNegative() {
		return Price{}, apperr.InvalidPriceErr
	}
	return Price{totalPrice: price}, nil
}
