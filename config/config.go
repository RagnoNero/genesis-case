package config

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func LoadEnvFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		os.Setenv(key, value)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading .env file: %v", err)
	}
}

func LoadConfig() *AppConfiguration {
	LoadEnvFile(".env")

	dbPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	dockerPort, _ := strconv.Atoi(os.Getenv("DOCKER_PORT"))
	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	return &AppConfiguration{
		DbConfig: DatabaseConfiguration{
			DbHost:     os.Getenv("DB_HOST"),
			DbPort:     dbPort,
			DbUser:     os.Getenv("DB_USER"),
			DbPassword: os.Getenv("DB_PASS"),
			DbName:     os.Getenv("DB_NAME"),
		},
		SmtpConfig: SmtpConfiguration{
			SmtpHost:     os.Getenv("SMTP_HOST"),
			SmtpPort:     smtpPort,
			SmtpEmail:    os.Getenv("SMTP_EMAIL"),
			SmtpUsername: os.Getenv("SMTP_USER"),
			SmtpPassword: os.Getenv("SMTP_PASS"),
		},
		WeatherApiKey: os.Getenv("WEATHER_API_KEY"),
		DockerPort:    dockerPort,
		AppUrl:        os.Getenv("APP_URL"),
	}
}
