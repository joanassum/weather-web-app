package adapters

import (
	"github.com/joanassum/weather-web-app/models"
)

type WeatherAPI interface {
	// GetWeatherInfo gets the weather information of a given city
	GetWeatherInfo(city string) (models.DomainWeatherInfo, error)
}
