package infrastructure

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joanassum/weather-web-app/models"
	"github.com/joho/godotenv"
)

const OPEN_API_URL string = "http://api.openweathermap.org/data/2.5/weather"

type openAPIWeatherInfo struct {
	Weather []openAPIWeather `json:"weather"`
	Main    openAPIMain      `json:"main"`
}

type openAPIWeather struct {
	Main     string `json:"main"`
	IconCode string `json:"icon"`
}

type openAPIMain struct {
	Temp     float64 `json:"temp"`
	Humidity int     `json:"humidity"`
}

func (w *openAPIWeatherInfo) toDomain() *models.DomainWeatherInfo {
	return &models.DomainWeatherInfo{
		Description: w.Weather[0].Main,
		Temp:        w.Main.Temp,
		Humidity:    w.Main.Humidity,
		IconURL:     "http://openweathermap.org/img/w/" + w.Weather[0].IconCode + ".png",
	}
}

type openWeatherAPI struct{}

// UseOpenWeatherAPI get a openWeatherAPI struct.
var UseOpenWeatherAPI openWeatherAPI

// GetWeatherInfo gets the weather information of a given city
func (openWeatherAPI) GetWeatherInfo(city string) (*models.DomainWeatherInfo, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return nil, err
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", OPEN_API_URL, nil)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	q := req.URL.Query()
	q.Add("appid", os.Getenv("OPEN_API_KEY"))
	q.Add("q", city)
	q.Add("units", "metric")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	fmt.Println(req.URL.String())

	defer resp.Body.Close()

	info := &openAPIWeatherInfo{}
	err = json.NewDecoder(resp.Body).Decode(info)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return info.toDomain(), nil
}
