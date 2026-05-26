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
	return NewPriceFromDecimal(price)
}

func NewPriceFromDecimal(price decimal.Decimal) (Price, error) {
	if price.IsZero() || price.IsNegative() {
		return Price{}, apperr.InvalidPriceErr
	}

	return Price{totalPrice: price}, nil
}

func (p Price) Decimal() decimal.Decimal {
	return p.totalPrice
}

func (p Price) String() string {
	return p.totalPrice.String()
}
