package app

import (
	"time"

	"github.com/Rhymond/go-money"
)

type Exchanger interface {
	Exchange(amount *money.Money, to *money.Currency) (*money.Money, error)
	ExchangeToDate(amount *money.Money, to *money.Currency, date time.Time) (*money.Money, error)
	UpdateRates() error
}
