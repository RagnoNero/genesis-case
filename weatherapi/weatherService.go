package weatherapi

import (
	"encoding/json"
	"fmt"
	"io"
	"weather-subscription/http"
	"weather-subscription/models"
)

const (
	baseWeatherURL  = "https://api.weatherapi.com/v1"
	currentEndpoint = "/current.json"
	queryParam      = "?q=%s&key=%s"
)

type WeatherService struct {
	client http.IHttpClient
	apiKey string
}

func NewWeatherService(client http.IHttpClient, apiKey string) *WeatherService {
	return &WeatherService{
		client: client,
		apiKey: apiKey,
	}
}

func (w *WeatherService) GetWeather(city string) (models.Weather, error) {
	url := fmt.Sprintf("%s%s%s", baseWeatherURL, currentEndpoint, fmt.Sprintf(queryParam, city, w.apiKey))
	resp, err := w.client.Get(url, nil)
	if err != nil {
		return models.Weather{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return models.Weather{}, fmt.Errorf("weather API error: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Weather{}, err
	}

	weather, err := parseWeather(body)
	if err != nil {
		return models.Weather{}, err
	}

	return *weather, nil
}

func parseWeather(apiResp []byte) (*models.Weather, error) {
	var data struct {
		Current struct {
			TempC     float64 `json:"temp_c"`
			Humidity  int     `json:"humidity"`
			Condition struct {
				Text string `json:"text"`
			} `json:"condition"`
		} `json:"current"`
	}

	err := json.Unmarshal(apiResp, &data)
	if err != nil {
		return nil, err
	}

	weather := &models.Weather{
		Temperature: int(data.Current.TempC),
		Humidity:    data.Current.Humidity,
		Description: data.Current.Condition.Text,
	}
	return weather, nil
}
