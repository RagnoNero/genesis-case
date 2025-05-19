package email

import "weather-subscription/models"

type IEmailSender interface {
	SendConfirmation(to string, token string) error
	SendWeather(to, city string, weather models.Weather, token string) error
}
