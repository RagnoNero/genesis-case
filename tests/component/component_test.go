package component

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"weather-subscription/api"
	"weather-subscription/config"
	"weather-subscription/models"
	"weather-subscription/scheduler"
	"weather-subscription/sql"
	"weather-subscription/sql/dto"
	"weather-subscription/tests/mocks"
	"weather-subscription/token"
	"weather-subscription/weatherapi"
)

func setupTestServer() (*api.Server, *http.ServeMux) {
	cfg := config.AppConfiguration{}

	repo := &mocks.DummySubRepo{
		Subs: []dto.SubscriptionDto{
			{
				Email:        "test@example.com",
				City:         "Kyiv",
				Frequency:    dto.Daily,
				Token:        "valid-token",
				Confirmed:    true,
				SubscribedAt: time.Now().Add(-2 * time.Hour),
				LastSentAt:   time.Now().Add(-time.Hour),
			},
		},
	}

	cache := scheduler.NewInMemorySubscriptionCache(repo)
	cache.ReloadCache()
	tokenGen := token.NewTokenGenerator(10)
	client := mocks.NewDummyHttp(ReturnCorrectResponse())
	weatherClient := weatherapi.NewWeatherService(client, "dummy")
	dbRepo := sql.NewDbRepository(repo)
	emailSender := &mocks.DummyEmailClient{}
	dummyScheduler := mocks.NewDummyScheduler()

	server := api.NewServer(cfg, *dbRepo, *weatherClient, tokenGen, emailSender, dummyScheduler)

	mux := http.NewServeMux()
	server.RegisterRoutesWithMux(mux)
	return server, mux
}

func TestSubscribeHandler(t *testing.T) {
	_, mux := setupTestServer()

	payload := `{"email":"user@mail.com","city":"Kyiv","frequency":"Daily"}`
	req := httptest.NewRequest(http.MethodPost, "/subscribe", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
}

func TestConfirmHandler(t *testing.T) {
	_, mux := setupTestServer()

	req := httptest.NewRequest(http.MethodGet, "/confirm/valid-token", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
}

func TestUnsubscribeHandler(t *testing.T) {
	_, mux := setupTestServer()

	req := httptest.NewRequest(http.MethodGet, "/unsubscribe/valid-token", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
}

func TestWeatherHandler(t *testing.T) {
	_, mux := setupTestServer()

	req := httptest.NewRequest(http.MethodGet, "/weather?city=Kyiv", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var w models.Weather
	err := json.NewDecoder(rr.Body).Decode(&w)
	if err != nil {
		t.Fatalf("failed to decode weather: %v", err)
	}

	if w.Temperature == 0 && w.Humidity == 0 && w.Description == "" {
		t.Error("expected non-empty weather data")
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
