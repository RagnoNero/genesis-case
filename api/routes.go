package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"weather-subscription/models"
)

func (s *Server) GetWeatherHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, "Missing 'city' query parameter", http.StatusBadRequest)
		return
	}

	weather, err := s.WeatherService.GetWeather(city)
	if err != nil {
		http.Error(w, "City not found or API error", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(weather)
	if err != nil {
		http.Error(w, "Failed to encode weather data", http.StatusInternalServerError)
		return
	}
}

func (s *Server) SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")

	subscription, err := Parse(contentType, r)

	if err != nil {
		http.Error(w, "Request is not valid", http.StatusBadRequest)
		return
	}

	if s.SqlClient.Subscription.IsSubscribed(subscription.Email) {
		http.Error(w, "Already subscribed", http.StatusConflict)
		return
	}

	_, err = s.WeatherService.GetWeather(subscription.City)
	if err != nil {
		http.Error(w, "City not found or API error", http.StatusNotFound)
		return
	}

	token, err := s.TokenGenerator.Generate(subscription.Email)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
	}

	subDto := subscription.ToDto()

	subDto.Token = token

	err = s.SqlClient.Subscription.CreateSubscription(subDto)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to subscribe: %v", err), http.StatusBadRequest)
		return
	}

	err = s.EmailService.SendConfirmation(subDto.Email, token)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to send email: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Subscription successful for %s\n", subDto.Email)
}

func (s *Server) ConfirmHandler(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimPrefix(r.URL.Path, "/confirm/")
	if token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	err := s.SqlClient.Subscription.ConfirmToken(token)
	if err != nil {
		http.Error(w, "Token not found", http.StatusNotFound)
		return
	}

	s.Scheduler.GetCache().AddToCache(token)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Subscription confirmed")
}

func (s *Server) UnsubscribeHandler(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimPrefix(r.URL.Path, "/unsubscribe/")
	if token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	err := s.SqlClient.Subscription.Unsubscribe(token)
	if err != nil {
		http.Error(w, "Token not found", http.StatusNotFound)
		return
	}

	s.Scheduler.GetCache().RemoveFromCache(token)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Unsubscribed successfully")
}

func (s *Server) RegisterRoutes() {
	http.HandleFunc("/weather", s.GetWeatherHandler)
	http.HandleFunc("/subscribe", s.SubscribeHandler)
	http.HandleFunc("/confirm/", s.ConfirmHandler)
	http.HandleFunc("/unsubscribe/", s.UnsubscribeHandler)
	http.Handle("/", http.FileServer(http.Dir("./static")))
}

func (s *Server) RegisterRoutesWithMux(mux *http.ServeMux) {
	mux.HandleFunc("/weather", s.GetWeatherHandler)
	mux.HandleFunc("/subscribe", s.SubscribeHandler)
	mux.HandleFunc("/confirm/", s.ConfirmHandler)
	mux.HandleFunc("/unsubscribe/", s.UnsubscribeHandler)
}

func (s *Server) StartServer(port string) {
	s.RegisterRoutes()
	log.Printf("Server listening on :%s...\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func Parse(contentType string, r *http.Request) (models.Subscription, error) {
	var email, city, freq string

	switch {
	case contentType == "application/json" || strings.HasPrefix(contentType, "application/json;"):
		var data struct {
			Email     string `json:"email"`
			City      string `json:"city"`
			Frequency string `json:"frequency"`
		}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			return models.Subscription{}, fmt.Errorf("invalid JSON: %w", err)
		}
		email, city, freq = data.Email, data.City, data.Frequency

	case contentType == "application/x-www-form-urlencoded" || strings.HasPrefix(contentType, "application/x-www-form-urlencoded;"):
		err := r.ParseForm()
		if err != nil {
			return models.Subscription{}, fmt.Errorf("invalid form data: %w", err)
		}
		email = r.FormValue("email")
		city = r.FormValue("city")
		freq = r.FormValue("frequency")

	case strings.HasPrefix(contentType, "multipart/form-data"):
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			return models.Subscription{}, fmt.Errorf("invalid multipart form data: %w", err)
		}
		email = r.FormValue("email")
		city = r.FormValue("city")
		freq = r.FormValue("frequency")

	default:
		return models.Subscription{}, fmt.Errorf("unsupported content-type: %s", contentType)
	}

	frequency, err := models.ParseFrequency(freq)
	if err != nil {
		return models.Subscription{}, fmt.Errorf("invalid frequency")
	}

	if email == "" || city == "" || freq == "" {
		return models.Subscription{}, fmt.Errorf("missing required fields")
	}

	return models.Subscription{
		Email:     email,
		City:      city,
		Frequency: frequency,
	}, nil
}
