package di

import (
	"weather-subscription/api"
	"weather-subscription/config"
	"weather-subscription/email"
	"weather-subscription/http"
	"weather-subscription/scheduler"
	"weather-subscription/sql"
	"weather-subscription/token"
	"weather-subscription/weatherapi"
)

type Application struct {
	AppConfiguration       config.AppConfiguration
	HttpClient             http.IHttpClient
	WeatherService         weatherapi.WeatherService
	SqlClient              sql.ISqlClient
	SubscriptionRepository sql.ISubscriptionRepository
	DbRepository           sql.DbRepository
	EmailSender            email.IEmailSender
	Server                 api.Server
	TokenGenerator         token.ITokenGenerator
	DynamicScheduler       scheduler.IScheduler
	InMemoryCache          scheduler.ISubscriptionCache
}

func NewApp(config config.AppConfiguration,
	httpClient http.IHttpClient,
	weatherservice weatherapi.WeatherService,
	sqlClient sql.ISqlClient,
	subscriptionRepo sql.ISubscriptionRepository,
	dbRepo sql.DbRepository,
	emailSender email.IEmailSender,
	server api.Server,
	tokenGenerator token.ITokenGenerator,
	scheduler scheduler.IScheduler,
	cache scheduler.ISubscriptionCache) *Application {

	return &Application{
		AppConfiguration:       config,
		HttpClient:             httpClient,
		WeatherService:         weatherservice,
		SqlClient:              sqlClient,
		SubscriptionRepository: subscriptionRepo,
		DbRepository:           dbRepo,
		EmailSender:            emailSender,
		Server:                 server,
		TokenGenerator:         tokenGenerator,
		DynamicScheduler:       scheduler,
		InMemoryCache:          cache,
	}
}
