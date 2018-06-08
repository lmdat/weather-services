package weatherapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Provider cho apixu.com
type ApiXuProvider struct {
	APIKey string
	URL    string
}

type ApiXuData struct {
	Current struct {
		CelsiusTemp float64 `json:"temp_c"`
	} `json:"current"`
}

// Implement hàm GetTemperature của WeatherProvider Interface
func (p ApiXuProvider) GetTemperature(city string) (float64, error) {
	res, err := http.Get(p.URL + p.APIKey + "&q=" + city)

	if err != nil || res.StatusCode != 200 {
		return 0, err
	}

	defer res.Body.Close()

	data := ApiXuData{}

	err = json.NewDecoder(res.Body).Decode(&data)

	if err != nil {
		return 0, err
	}

	fmt.Println("apixu: ", data.Current.CelsiusTemp)
	return data.Current.CelsiusTemp, err
}
