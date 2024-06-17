package dto

type Weather struct {
	Current WeatherCurrent `json:"current"`
}

type WeatherCurrent struct {
	LastUpdated string  `json:"last_updated"`
	TempC       float64 `json:"temp_c"`
	TempF       float64 `json:"temp_f"`
}
