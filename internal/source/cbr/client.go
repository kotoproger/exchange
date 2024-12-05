package cbr

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Answer struct {
	Rates map[string]Rate `json:"Valute"` //nolint:tagliatelle
}

type Rate struct {
	CurrencyCode string  `json:"CharCode"` //nolint:tagliatelle
	Nominal      int32   `json:"Nominal"`  //nolint:tagliatelle
	Rate         float64 `json:"Value"`    //nolint:tagliatelle
}

func getRates(out chan<- Rate, url string) {
	defer close(out)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Ошибка запроса", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Ошибка HTTP-ответа: %d\n", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка чтения", err)
		return
	}

	response := Answer{}

	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Ошибка десериализации", err)
		return
	}

	for _, rate := range response.Rates {
		out <- rate
	}
}
