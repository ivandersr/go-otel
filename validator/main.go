package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/ivandersr/validator/pkg/utils"
)

type WeatherRequest struct {
	Cep string `json:"cep"`
}

type WeatherResponse struct {
	TempC string `json:"temp_C"`
	TempF string `json:"temp_F"`
	TempK string `json:"temp_K"`
}

const (
	unprocessableError = "invalid zipcode"
	unexpectedError    = "there was an error trying to query for current temperature"
	notFoundError      = "can not find zipcode"
)

func main() {
	http.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var reqBody WeatherRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		parsedCep, err := utils.ValidateCep(reqBody.Cep)
		if err != nil {
			http.Error(w, unprocessableError, http.StatusUnprocessableEntity)
			return
		}

		weatherURL := os.Getenv("WEATHER_API")
		resp, err := http.Get(weatherURL + parsedCep)
		if err != nil {
			http.Error(w, unexpectedError, http.StatusInternalServerError)
			return
		}

		if resp.StatusCode != http.StatusOK {
			if resp.StatusCode == http.StatusNotFound {
				http.Error(w, notFoundError, http.StatusNotFound)
				return
			}
			http.Error(w, unexpectedError, http.StatusInternalServerError)
		}

		var weather WeatherResponse
		err = json.NewDecoder(resp.Body).Decode(&weather)
		if err != nil {
			http.Error(w, unexpectedError, http.StatusInternalServerError)
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

	http.ListenAndServe(":8080", nil)
}
