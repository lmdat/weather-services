package weatherapi

// Weather là interface để 3 service implement hàm GetTemperature
type WeatherProvider interface {
	GetTemperature(city string) (float64, error)
}
