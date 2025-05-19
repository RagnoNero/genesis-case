package unit

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"weather-subscription/tests/mocks"
	"weather-subscription/weatherapi"
)

func GetWeather_ReturnsBadRequest(t *testing.T) {
	httpClient := NewDummyHttp(ReturnBadRequest())
	weatherApi := weatherapi.NewWeatherService(httpClient, "dummy-key")

	_, err := weatherApi.GetWeather("Kyiv")

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetWeather_ReturnsCorrectWeather(t *testing.T) {
	client := NewDummyHttp(ReturnCorrectResponse())
	service := weatherapi.NewWeatherService(client, "dummy")

	weather, err := service.GetWeather("Kyiv")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if weather.Temperature != 22 {
		t.Errorf("expected Temperature 22, got %d", weather.Temperature)
	}
	if weather.Humidity != 60 {
		t.Errorf("expected Humidity 60, got %d", weather.Humidity)
	}
	if weather.Description != "Partly cloudy" {
		t.Errorf("expected Description 'Partly cloudy', got '%s'", weather.Description)
	}
}

func NewDummyHttp(retFunc func() (*http.Response, error)) *mocks.DummyHttpClient {
	return &mocks.DummyHttpClient{
		ReturnFunc: retFunc,
	}
}

func ReturnCorrectResponse() func() (*http.Response, error) {
	jsonResponse := `{
		"current": {
			"temp_c": 22.5,
			"humidity": 60,
			"condition": {
				"text": "Partly cloudy"
			}
		}
	}`

	return func() (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Body:       io.NopCloser(strings.NewReader(jsonResponse)),
		}, nil
	}
}

func ReturnBadRequest() func() (*http.Response, error) {
	return func() (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Status:     "400 Bad Request",
		}, fmt.Errorf("weather API error")
	}
}
