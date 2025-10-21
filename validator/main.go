package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/ivandersr/validator/internal/tracing"
	"github.com/ivandersr/validator/pkg/utils"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

type WeatherRequest struct {
	Cep string `json:"cep"`
}

type WeatherResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

const (
	unprocessableError = "invalid zipcode"
	unexpectedError    = "there was an error trying to query for current temperature"
	notFoundError      = "can not find zipcode"
)

func main() {
	collectorURL := os.Getenv("COLLECTOR_URL")
	tp := tracing.InitTracer("validator-service", collectorURL)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()

	mux := http.NewServeMux()
	mux.Handle("/weather", otelhttp.NewHandler(http.HandlerFunc(weatherHandler), "validator-handler"))

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("validator")
	ctx, validatorSpan := tracer.Start(r.Context(), "validator-span")
	defer validatorSpan.End()

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
	instrumentedClient := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, weatherURL+parsedCep, nil)
	resp, err := instrumentedClient.Do(req)
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
}
