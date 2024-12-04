package source

import money "github.com/Rhymond/go-money"

type ExchangeRate struct {
	From money.Currency
	To   money.Currency
	Rate float64
}

type ExchangeSource interface {
	Get() <-chan ExchangeRate
}
