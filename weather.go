package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const OPEN_API_URL string = "http://api.openweathermap.org/data/2.5/weather"

type openAPIWeatherInfo struct {
	Weather []openAPIWeather `json:"weather"`
	Main    OpenAPIMain      `json:"main"`
}

type openAPIWeather struct {
	Main     string `json:"main"`
	IconCode string `json:"icon"`
}

type OpenAPIMain struct {
	Temp     float64 `json:"temp"`
	Humidity int     `json:"humidity"`
}

type DomainWeatherInfo struct {
	Description string
	Temp        float64
	IconURL     string
	Humidity    int
}

func (w *openAPIWeatherInfo) ToDomain() *DomainWeatherInfo {
	return &DomainWeatherInfo{
		Description: w.Weather[0].Main,
		Temp:        w.Main.Temp,
		Humidity:    w.Main.Humidity,
		IconURL:     "http://openweathermap.org/img/w/" + w.Weather[0].IconCode + ".png",
	}
}

type WeatherAPI interface {
	getWeatherInfo(city string) (DomainWeatherInfo, error)
}

type OpenWeatherAPI struct{}

var UseOpenWeatherAPI OpenWeatherAPI

func (OpenWeatherAPI) getWeatherInfo(city string) (*DomainWeatherInfo, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", OPEN_API_URL, nil)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	q := req.URL.Query()
	q.Add("appid", OPEN_API_KEY)
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

	return info.ToDomain(), nil
}
