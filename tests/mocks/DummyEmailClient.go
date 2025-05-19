package mocks

import "weather-subscription/models"

type DummyEmailClient struct {
}

func (d *DummyEmailClient) SendConfirmation(to, token string) error {
	return nil
}

func (d *DummyEmailClient) SendWeather(to, city string, weather models.Weather, token string) error {
	return nil
}

func NewDummyEmailClient() *DummyEmailClient {
	return &DummyEmailClient{}
}
