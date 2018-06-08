package weatherapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Provider cho weatherbit.io
type WeatherBitProvider struct {
	APIKey string
	URL    string
}

type WeatherBitData struct {
	Current []struct {
		CelsiusTemp float64 `json:"temp"`
	} `json:"data"`
}

// Implement hàm GetTemperature của Weather Interface
func (p WeatherBitProvider) GetTemperature(city string) (float64, error) {
	res, err := http.Get(p.URL + p.APIKey + "&city=" + city)

	if err != nil || res.StatusCode != 200 {
		return 0, err
	}

	defer res.Body.Close()

	data := WeatherBitData{}

	err = json.NewDecoder(res.Body).Decode(&data)

	if err != nil {
		return 0, err
	}

	fmt.Println("weatherbit: ", data.Current[0].CelsiusTemp)
	return data.Current[0].CelsiusTemp, err
}
