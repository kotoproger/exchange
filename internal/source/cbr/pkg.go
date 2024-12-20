package cbr

import (
	"github.com/Rhymond/go-money"
	"github.com/kotoproger/exchange/internal/source"
)

type Cbr struct {
	url string
}

func NewCbr(url string) *Cbr {
	return &Cbr{url: url}
}

func (c Cbr) Get() <-chan source.ExchangeRate {
	output := make(chan source.ExchangeRate)
	input := make(chan Rate)
	go getRates(input, c.url)
	go tranformRates(output, input)

	return output
}

func tranformRates(in chan<- source.ExchangeRate, input <-chan Rate) {
	defer close(in)

	rubCurrency := money.GetCurrency("RUB")
	if rubCurrency == nil {
		return
	}
	rawRates := []Rate{}
	for rate := range input {
		rawRates = append(rawRates, rate)
	}
	for _, rate := range rawRates {
		toCurrency := money.GetCurrency(rate.CurrencyCode)
		if toCurrency == nil {
			continue
		}
		in <- source.ExchangeRate{
			From: *rubCurrency,
			To:   *toCurrency,
			Rate: float64(rate.Nominal) / rate.Rate,
		}
		in <- source.ExchangeRate{
			From: *toCurrency,
			To:   *rubCurrency,
			Rate: rate.Rate / float64(rate.Nominal),
		}
		for _, rate2 := range rawRates {
			if rate.CurrencyCode == rate2.CurrencyCode {
				continue
			}
			to2Currency := money.GetCurrency(rate2.CurrencyCode)
			in <- source.ExchangeRate{
				From: *to2Currency,
				To:   *toCurrency,
				Rate: (float64(rate.Nominal) / rate.Rate) / (float64(rate2.Nominal) / rate2.Rate),
			}
		}
	}
}
