package email

import (
	"fmt"
	"net/smtp"
	"strings"
	"weather-subscription/config"
	"weather-subscription/models"
)

type SmtpEmailSender struct {
	SmtpConfig config.SmtpConfiguration
	BaseUrl    string
}

func NewSmtpEmailSender(config config.SmtpConfiguration, baseUrl string) *SmtpEmailSender {
	return &SmtpEmailSender{
		SmtpConfig: config,
		BaseUrl:    baseUrl,
	}
}

func (s *SmtpEmailSender) SendConfirmation(to, token string) error {
	auth := smtp.PlainAuth("", s.SmtpConfig.SmtpUsername, s.SmtpConfig.SmtpPassword, s.SmtpConfig.SmtpHost)

	body := generateConfirmEmailBody(s.BaseUrl, token)

	msg := []byte(fmt.Sprintf(
		"To: %s\r\nSubject: %s\r\nMIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n%s",
		to, "Subscription Confirmation", body,
	))

	addr := fmt.Sprintf("%s:%d", s.SmtpConfig.SmtpHost, s.SmtpConfig.SmtpPort)
	recipients := strings.Split(to, ",")

	err := smtp.SendMail(addr, auth, s.SmtpConfig.SmtpEmail, recipients, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (s *SmtpEmailSender) SendWeather(to, city string, weather models.Weather, token string) error {
	auth := smtp.PlainAuth("", s.SmtpConfig.SmtpUsername, s.SmtpConfig.SmtpPassword, s.SmtpConfig.SmtpHost)

	body := generateWeatherEmailBody(s.BaseUrl, city, weather, token)

	msg := []byte(fmt.Sprintf(
		"To: %s\r\nSubject: %s\r\nMIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n%s",
		to, "Weather report", body,
	))

	addr := fmt.Sprintf("%s:%d", s.SmtpConfig.SmtpHost, s.SmtpConfig.SmtpPort)
	recipients := strings.Split(to, ",")

	err := smtp.SendMail(addr, auth, s.SmtpConfig.SmtpEmail, recipients, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
