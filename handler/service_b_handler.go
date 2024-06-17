package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/leoseiji/go-tracing/dto"
	"go.opentelemetry.io/otel"
)

var ErrCEPNotFound = fmt.Errorf("can not find zipcode")
var ErrCEPInvalid = fmt.Errorf("invalid zipcode")

func GetWeatherHandler(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("weather-service-b")
	_, span := tracer.Start(r.Context(), "GetWeatherHandler")
	defer span.End()

	cep := r.PathValue("cep")

	if !isCepValid(cep) {
		fmt.Printf("CEP %s is invalid", cep)
		http.Error(w, ErrCEPInvalid.Error(), http.StatusUnprocessableEntity)
		return
	}

	location, err := getLocationByCEP(cep)
	if errors.Is(err, ErrCEPNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	weather, err := getWeatherByLocation(location.Location)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	weatherResponse := dto.NewCEPWeatherResponse(location, weather)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weatherResponse)
}

func isCepValid(cep string) bool {
	if cep == "" {
		return false
	}
	if len(cep) != 8 {
		return false
	}
	if !regexp.MustCompile(`^[0-9]*$`).MatchString(cep) {
		return false
	}
	fmt.Printf("CEP %s is valid", cep)
	return true
}

func getLocationByCEP(cep string) (*dto.Location, error) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("error creating ViaCEP request. Err:%s", err.Error())
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("error executing ViaCEP request. Err:%s", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {

	case http.StatusOK:
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("error while reading ViaCEP result. Err:%s", err.Error())
			return nil, err
		}

		var location *dto.Location
		if err = json.Unmarshal(body, &location); err != nil {
			log.Printf("error while converting ViaCEP result. Err:%s", err.Error())
			return nil, err
		}
		if location.CEP == "" {
			return nil, ErrCEPNotFound
		}
		return location, nil

	case http.StatusNotFound:
		return nil, ErrCEPNotFound

	default:
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

}

func getWeatherByLocation(location string) (*dto.Weather, error) {
	location = strings.Replace(location, " ", "%20", -1)
	reqUrl := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=e6c189ac26084b8a84213356241706&q=%s", url.PathEscape(location))

	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Printf("error creating weatherAPI request. Err:%s", err.Error())
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("error executing weatherAPI request. Err:%s", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("error while getting weatherAPI result. Status: %s, Body: %s", resp.Status, string(body))

		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error while reading weatherAPI result. Err:%s", err.Error())
		return nil, err
	}

	var weather *dto.Weather
	if err = json.Unmarshal(body, &weather); err != nil {
		log.Printf("error while converting weatherAPI result. Err:%s", err.Error())
		return nil, err
	}
	return weather, nil
}
