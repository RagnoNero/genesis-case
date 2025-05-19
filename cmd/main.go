package main

import (
	"fmt"
	"log"
	"time"

	"weather-subscription/api"
	"weather-subscription/config"
	"weather-subscription/email"
	"weather-subscription/http"
	"weather-subscription/scheduler"
	"weather-subscription/sql"
	"weather-subscription/token"
	"weather-subscription/weatherapi"
)

func main() {
	configuration := config.LoadConfig()
	dbClient, err := sql.NewPostgresClient(configuration.DbConfig)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	baseUrl := fmt.Sprintf("%s:%d", configuration.AppUrl, configuration.DockerPort)

	subscriptionRepo := sql.NewSubscriptionRepository(dbClient.GetDb())
	dbRepo := sql.NewDbRepository(subscriptionRepo)
	httpClient := http.NewHttpClient(10 * time.Second)
	weatherService := weatherapi.NewWeatherService(httpClient, configuration.WeatherApiKey)
	tokenService := token.NewTokenGenerator(10)
	mailService := email.NewSmtpEmailSender(configuration.SmtpConfig, baseUrl)
	memoryCache := scheduler.NewInMemorySubscriptionCache(dbRepo.Subscription)
	scheduler := scheduler.NewDynamicScheduler(*dbRepo, mailService, *weatherService, memoryCache)

	s := api.NewServer(
		*configuration,
		*dbRepo,
		*weatherService,
		tokenService,
		mailService,
		scheduler,
	)

	s.Scheduler.Start()
	log.Println("Server started at :8080")
	s.StartServer("8080")
}
