package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/ivandersr/go-weather/internal/tracing"
	"github.com/ivandersr/go-weather/pkg/cep"
	"github.com/ivandersr/go-weather/pkg/weather"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

func main() {
	collectorURL := os.Getenv("COLLECTOR_URL")
	tp := tracing.InitTracer("go-weather-service", collectorURL)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()

	mux := http.NewServeMux()

	mux.Handle("/weather", otelhttp.NewHandler(http.HandlerFunc(weatherHandler), "weather-handler"))

	if err := http.ListenAndServe(":8081", mux); err != nil {
		panic(err)
	}
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("go-weather")
	_, cepSpan := tracer.Start(context.Background(), "cep-span")

	requestedCep := r.URL.Query().Get("cep")
	cepURL := os.Getenv("CEP_API")
	result, err := cep.GetLocation(cepURL, requestedCep)
	if err != nil {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}
	cepSpan.End()

	_, weatherSpan := tracer.Start(context.Background(), "weather-span")
	defer weatherSpan.End()
	weatherURL := os.Getenv("WEATHER_API")
	weather, err := weather.GetWeather(weatherURL, result.City, result.State, result.CityOriginal)
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
}
