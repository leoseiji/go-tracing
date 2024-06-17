package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/leoseiji/go-tracing/dto"
	"go.opentelemetry.io/otel"
)

var ErrInternalServerError = fmt.Errorf("internal server error")

func PostWeatherHandler(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("weather-service-a")
	ctx, span := tracer.Start(r.Context(), "PostWeatherHandler")
	defer span.End()

	var weatherCepRequest dto.WeatherCepRequest
	if err := json.NewDecoder(r.Body).Decode(&weatherCepRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !isCepValid(weatherCepRequest.Cep) {
		fmt.Printf("CEP %s is invalid", weatherCepRequest.Cep)
		http.Error(w, ErrCEPInvalid.Error(), http.StatusUnprocessableEntity)
		return
	}

	url := fmt.Sprintf("http://localhost:8080/weather-service-b/%s", weatherCepRequest.Cep)
	cepWeatherReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Printf("error while creating request: %s", err)
		http.Error(w, ErrInternalServerError.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := http.DefaultClient.Do(cepWeatherReq)
	if err != nil {
		log.Printf("error while making request: %s", err)
		http.Error(w, ErrInternalServerError.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, ErrInternalServerError.Error(), http.StatusInternalServerError)
			return
		}
		var location *dto.CEPWeatherResponse
		if err = json.Unmarshal(body, &location); err != nil {
			log.Printf("error while unmarshaling response: %s", err)
			http.Error(w, ErrInternalServerError.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(body)
		return

	case http.StatusNotFound:
		log.Printf("error while making request: %s", err)
		http.Error(w, ErrCEPNotFound.Error(), http.StatusNotFound)
		return

	default:
		log.Printf("unexpected error: %s", err)
		http.Error(w, ErrInternalServerError.Error(), http.StatusInternalServerError)
		return
	}

}
