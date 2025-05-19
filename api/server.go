package api

import (
	"fmt"
	"net"
	"weather-subscription/config"
	"weather-subscription/email"
	"weather-subscription/scheduler"
	"weather-subscription/sql"
	"weather-subscription/token"
	"weather-subscription/weatherapi"
)

type Server struct {
	Configuration  config.AppConfiguration
	SqlClient      sql.DbRepository
	WeatherService weatherapi.WeatherService
	TokenGenerator token.ITokenGenerator
	EmailService   email.IEmailSender
	Scheduler      scheduler.IScheduler
}

func NewServer(config config.AppConfiguration,
	sqlClient sql.DbRepository,
	weather weatherapi.WeatherService,
	token token.ITokenGenerator,
	email email.IEmailSender,
	scheduler scheduler.IScheduler) *Server {
	return &Server{
		Configuration:  config,
		SqlClient:      sqlClient,
		WeatherService: weather,
		TokenGenerator: token,
		EmailService:   email,
		Scheduler:      scheduler,
	}
}

func GetLocalIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if !ok || ipNet.IP.IsLoopback() {
			continue
		}
		ip := ipNet.IP.To4()
		if ip == nil {
			continue
		}
		return ip.String(), nil
	}

	return "", fmt.Errorf("no IP address found")
}
