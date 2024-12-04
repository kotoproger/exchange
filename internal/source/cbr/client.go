package cbr

type Answer struct {
	Rates map[string]Rate `json:"Valute"` //nolint:tagliatelle
}

type Rate struct {
	CurrencyCode string  `json:"CharCode"` //nolint:tagliatelle
	Nominal      int32   `json:"Nominal"`  //nolint:tagliatelle
	Rate         float64 `json:"Value"`    //nolint:tagliatelle
}

func getRates() []Rate {
	return []Rate{}
}
