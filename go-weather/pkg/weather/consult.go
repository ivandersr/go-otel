package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ivandersr/go-weather/pkg/utils"
)

type Weather struct {
	Current struct {
		TempC float64 `json:"temp_c"`
		TempF float64 `json:"temp_f"`
	} `json:"current"`
}

type WeatherResult struct {
	TempC string `json:"temp_C"`
	TempF string `json:"temp_F"`
	TempK string `json:"temp_K"`
}

func GetWeather(weatherURL, city, state string) (*WeatherResult, error) {
	queryString := fmt.Sprintf("%s,%s", city, state)
	url := utils.InjectArg(queryString, weatherURL)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var w Weather
	err = json.Unmarshal(body, &w)
	if err != nil {
		return nil, err
	}

	return &WeatherResult{
		TempC: fmt.Sprintf("%.2f", w.Current.TempC),
		TempF: fmt.Sprintf("%.2f", w.Current.TempF),
		TempK: fmt.Sprintf("%.2f", w.Current.TempC+273.15),
	}, nil
}
