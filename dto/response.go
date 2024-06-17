package dto

type CEPWeatherResponse struct {
	Location                string  `json:"city"`
	TemperatureInCelcius    float64 `json:"temp_C"`
	TemperatureInFahrenheit float64 `json:"temp_F"`
	TemperatureInKelvin     float64 `json:"temp_K"`
}

func NewCEPWeatherResponse(location *Location, weather *Weather) *CEPWeatherResponse {
	return &CEPWeatherResponse{
		Location:                location.Location,
		TemperatureInCelcius:    weather.Current.TempC,
		TemperatureInFahrenheit: weather.Current.TempF,
		TemperatureInKelvin:     weather.Current.TempC + 273.15,
	}
}
