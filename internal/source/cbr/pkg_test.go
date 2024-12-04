package cbr

import (
	"testing"

	"github.com/Rhymond/go-money"
	"github.com/kotoproger/exchange/internal/source"
	"github.com/stretchr/testify/assert"
)

func TestTranformRates(t *testing.T) {
	testCases := []struct {
		name   string
		input  []Rate
		output []source.ExchangeRate
	}{
		{name: "empty", input: []Rate{}, output: []source.ExchangeRate{}},
		{
			name:  "one currency",
			input: []Rate{{CurrencyCode: "USD", Nominal: 2, Rate: 50}},
			output: []source.ExchangeRate{
				{
					From: *money.GetCurrency("RUB"),
					To:   *money.GetCurrency("USD"),
					Rate: 0.04,
				},
				{
					From: *money.GetCurrency("USD"),
					To:   *money.GetCurrency("RUB"),
					Rate: 25,
				},
			},
		},
		{
			name: "two currencies",
			input: []Rate{
				{CurrencyCode: "USD", Nominal: 1, Rate: 104.2361},
				{CurrencyCode: "KZT", Nominal: 100, Rate: 19.9},
			},
			output: []source.ExchangeRate{
				{
					From: *money.GetCurrency("RUB"),
					To:   *money.GetCurrency("USD"),
					Rate: float64(1) / 104.2361,
				},
				{
					From: *money.GetCurrency("USD"),
					To:   *money.GetCurrency("RUB"),
					Rate: 104.2361,
				},
				{
					From: *money.GetCurrency("KZT"),
					To:   *money.GetCurrency("USD"),
					Rate: (float64(1) / 104.2361) / (float64(100) / 19.9),
				},
				{
					From: *money.GetCurrency("RUB"),
					To:   *money.GetCurrency("KZT"),
					Rate: float64(100) / 19.9,
				},
				{
					From: *money.GetCurrency("KZT"),
					To:   *money.GetCurrency("RUB"),
					Rate: 19.9 / float64(100),
				},
				{
					From: *money.GetCurrency("USD"),
					To:   *money.GetCurrency("KZT"),
					Rate: (float64(100) / 19.9) / (float64(1) / 104.2361),
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			output := make(chan source.ExchangeRate, len(testCase.output))
			tranformRates(output, testCase.input)
			actualRates := make([]source.ExchangeRate, 0)
			for rate := range output {
				actualRates = append(actualRates, rate)
			}

			assert.Equal(t, testCase.output, actualRates)
		})
	}
}
