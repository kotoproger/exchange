package cbr

type CbrAnswer struct {
	Rates map[string]CbrRate `json:"Valute"`
}

type CbrRate struct {
	CurrencyCode string  `json:"CharCode"`
	Nominal      int32   `json:"Nominal"`
	Rate         float64 `json:"Value"`
}

func getRates() []CbrRate {
	return []CbrRate{}
}
