package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	wea "weather-services/weatherapi"

	"github.com/gorilla/mux"
)

// Danh sách các Provider
type ProviderList []wea.WeatherProvider

type TemperatureData struct {
	CityName       string  `json:"city_name"`
	CelsiusTemp    float64 `json:"celsius_temp"`
	KelvinTemp     float64 `json:"kelvin_temp"`
	FahrenheitTemp float64 `json:"fahrenheit_temp"`
}

// Lấy dữ liệu nhiệt độ và tính trung bình
func (list ProviderList) temperature(city string) float64 {

	// Tạo channel để hứng data và error trả về từ routine
	chanTemp := make(chan float64)
	chanErr := make(chan error)

	// Tạo các routine để thực hiện việc lấy data nhiệt độ từ 3 nguồn:
	// -Open Weather Map
	// -ApiXu
	// -Weather Bit
	for _, p := range list {
		// Run routine
		go func(w wea.WeatherProvider) {
			temp, err := w.GetTemperature(city)
			if err != nil {
				chanErr <- err
				return
			}
			// Đẩy dữ liệu nhiệt độ vào channel
			chanTemp <- temp
		}(p)
	}

	total := 0.0
	k := 0
	// Lấy dữ liệu nhiệt độ từ các channel (nếu có)
	for i := 0; i < len(list); i++ {
		select {
		case temp := <-chanTemp:
			if temp > 0 {
				total += temp
				k++
			}

		case err := <-chanErr:
			panic(err)
		}

	}
	// Sau đó tính trung bình nhiệt độ và trả kết quả
	return total / float64(k)
}

func main() {

	// Tạo provider để gọi api openweathermap.org
	openWeatherMap := wea.OpenWeatherMapProvider{
		APIKey: "b1668a59088cb0267b3cf221325408f7",
		URL:    "https://api.openweathermap.org/data/2.5/weather?appid=",
	}

	// Tạo provider để gọi api apixu.com
	apiXu := wea.ApiXuProvider{
		APIKey: "9da59343499f423cb3d90407180706",
		URL:    "https://api.apixu.com/v1/current.json?key=",
	}

	// Tạo provider để gọi api weatherbit.io
	weatherBit := wea.WeatherBitProvider{
		APIKey: "8da0ecc3363e4b51a39e4827fbcc71c5",
		URL:    "https://api.weatherbit.io/v2.0/current?key=",
	}

	// Danh sách chứa các service
	providerList := ProviderList{
		openWeatherMap,
		apiXu,
		weatherBit,
	}

	// cites := map[string]string{
	// 	"saigon":    "Ho Chi Minh",
	// 	"hanoi":     "Hanoi",
	// 	"hue":       "Hue",
	// 	"da-nang":   "Da Nang",
	// 	"nha-trang": "Nha Trang",
	// }

	// for key, val := range cites {
	// 	temp := providerList.temperature(key)
	// 	fmt.Printf("Temperature of %s is %f\n", val, temp)
	// }

	// Xử lý Rest API sử dụng thư viện Gorilla Mux
	r := mux.NewRouter()
	r.HandleFunc("/api/temperature/{city}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		city := vars["city"]

		// Lấy nhiệt độ
		tempC := providerList.temperature(city)
		tempK := tempC + 273.15
		tempF := (tempC * 1.8) + 32

		data := TemperatureData{
			CityName:       city,
			CelsiusTemp:    tempC,
			KelvinTemp:     tempK,
			FahrenheitTemp: tempF,
		}
		fmt.Printf("Temperature of %s is %f Celsius, %f Kelvin, %f Fahrenheit\n\n", city, tempC, tempK, tempF)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)

	}).Methods("GET")

	port := 9000
	fmt.Printf("Server is listening at port: %d\n", port)
	log.Fatal(http.ListenAndServe(":"+fmt.Sprint(port), r))
}
