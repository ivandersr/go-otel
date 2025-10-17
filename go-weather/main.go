package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/ivandersr/go-weather/pkg/cep"
	"github.com/ivandersr/go-weather/pkg/weather"
)

func main() {
	http.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
		requestedCep := r.URL.Query().Get("cep")

		cepURL := os.Getenv("CEP_API")
		result, err := cep.GetLocation(cepURL, requestedCep)
		if err != nil {
			http.Error(w, "can not find zipcode", http.StatusNotFound)
			return
		}

		weatherURL := os.Getenv("WEATHER_API")
		weather, err := weather.GetWeather(weatherURL, result.City, result.State)
		if err != nil {
			http.Error(w, "can not find weather for the given zipcode", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(weather)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.ListenAndServe(":8081", nil)
}
