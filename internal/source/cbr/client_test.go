package cbr

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRatesSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{
"Date": "2024-12-06T11:30:00+03:00",
"PreviousDate": "2024-12-05T11:30:00+03:00",
"PreviousURL": "//www.cbr-xml-daily.ru/archive/2024/12/05/daily_json.js",
"Timestamp": "2024-12-05T20:00:00+03:00",
"Valute": {
"AUD": {
"ID": "R01010",
"NumCode": "036",
"CharCode": "AUD",
"Nominal": 1,
"Name": "Австралийский доллар",
"Value": 66.5067,
"Previous": 67.1072
},
"AZN": {
"ID": "R01020A",
"NumCode": "944",
"CharCode": "AZN",
"Nominal": 1,
"Name": "Азербайджанский манат",
"Value": 60.8139,
"Previous": 61.3154
}}}`)
	}))
	defer ts.Close()
	expected := []Rate{
		{CurrencyCode: "AUD", Nominal: 1, Rate: 66.5067},
		{CurrencyCode: "AZN", Nominal: 1, Rate: 60.8139},
	}

	input := make(chan Rate, 2)
	getRates(input, ts.URL)
	rates := make([]Rate, 0)
	for rate := range input {
		rates = append(rates, rate)
	}

	assert.ElementsMatch(t, expected, rates)
}

func TestGetRatesMalformedJson(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{
"Date": "2024-12-06T11:30:00+03:00",
"PreviousDate": "2024-12-05T11:30:00+03:00",
"PreviousURL": "//www.cbr-xml-daily.ru/archive/2024/12/05/daily_json.js",
"Timestamp": "2024-12-05T20:00:00+03:00",
"Valute": {
"AUD": {
"ID": "R01010",
"NumCode": "036",
"CharCode": "AUD",
"Nominal": 1,
"Name": "Австралийский доллар",
"Value": 66.5067,
"Previous": 67.1072
},
"AZN": {
"ID": "R01020A",
"NumCode": "944",
"CharCode": "AZN",
"Nominal": 1,
"Name": "Азербайджанский манат",
"Value": 60.8139,
"Previous": 61.3154
}}`)
	}))
	defer ts.Close()
	expected := []Rate{}

	input := make(chan Rate)
	getRates(input, ts.URL)
	rates := make([]Rate, 0)
	for rate := range input {
		rates = append(rates, rate)
	}

	assert.ElementsMatch(t, expected, rates)
}

func TestGetRatesErrorneus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		fmt.Fprintln(w, `{
"Date": "2024-12-06T11:30:00+03:00",
"PreviousDate": "2024-12-05T11:30:00+03:00",
"PreviousURL": "//www.cbr-xml-daily.ru/archive/2024/12/05/daily_json.js",
"Timestamp": "2024-12-05T20:00:00+03:00",
"Valute": {
"AUD": {
"ID": "R01010",
"NumCode": "036",
"CharCode": "AUD",
"Nominal": 1,
"Name": "Австралийский доллар",
"Value": 66.5067,
"Previous": 67.1072
},
"AZN": {
"ID": "R01020A",
"NumCode": "944",
"CharCode": "AZN",
"Nominal": 1,
"Name": "Азербайджанский манат",
"Value": 60.8139,
"Previous": 61.3154
}}}`)
	}))
	defer ts.Close()
	expected := []Rate{}

	input := make(chan Rate)
	getRates(input, ts.URL)
	rates := make([]Rate, 0)
	for rate := range input {
		rates = append(rates, rate)
	}

	assert.ElementsMatch(t, expected, rates)
}
