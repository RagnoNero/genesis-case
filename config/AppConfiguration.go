package config

type AppConfiguration struct {
	DbConfig      DatabaseConfiguration
	SmtpConfig    SmtpConfiguration
	WeatherApiKey string
	AppUrl        string
	DockerPort    int
}
