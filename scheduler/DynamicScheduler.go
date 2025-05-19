package scheduler

import (
	"log"
	"time"
	"weather-subscription/email"
	"weather-subscription/sql"
	"weather-subscription/sql/dto"
	"weather-subscription/weatherapi"
)

type DynamicScheduler struct {
	Repository    sql.DbRepository
	EmailSender   email.IEmailSender
	WeatherClient weatherapi.WeatherService
	MemoryCache   ISubscriptionCache
}

func NewDynamicScheduler(repo sql.DbRepository, sender email.IEmailSender, weatherClient weatherapi.WeatherService, cache ISubscriptionCache) *DynamicScheduler {
	return &DynamicScheduler{
		Repository:    repo,
		EmailSender:   sender,
		WeatherClient: weatherClient,
		MemoryCache:   cache,
	}
}

func (s *DynamicScheduler) Start() {
	if err := s.MemoryCache.ReloadCache(); err != nil {
		log.Println("initial cache load failed:", err)
	}

	go s.reloadLoop()
	go s.dispatchLoop()
}

func (s *DynamicScheduler) reloadLoop() {
	for {
		err := s.MemoryCache.ReloadCache()
		if err != nil {
			log.Println("failed to reload cache:", err)
		}
		time.Sleep(24 * time.Hour)
	}
}

func (s *DynamicScheduler) dispatchLoop() {
	for {
		subs := s.MemoryCache.GetAll()
		now := time.Now()

		for _, sub := range subs {
			var interval time.Duration
			if sub.Frequency == dto.Hourly {
				interval = time.Hour
			} else {
				interval = 24 * time.Hour
			}

			last := sub.LastSentAt
			if last.IsZero() {
				last = sub.SubscribedAt
			}

			if now.Sub(last) >= interval {
				go func(si dto.SubscriptionDto) {
					weather, err := s.WeatherClient.GetWeather(si.City)
					if err != nil {
						log.Println("weather error:", err)
						return
					}
					if err := s.EmailSender.SendWeather(si.Email, si.City, weather, si.Token); err == nil {
						_ = s.Repository.Subscription.UpdateLastSent(si.Email)
						s.MemoryCache.UpdateLastSent(si.Email, time.Now())
					}
				}(sub)
			}
		}

		time.Sleep(1 * time.Minute)
	}
}

func (s *DynamicScheduler) GetCache() ISubscriptionCache {
	return s.MemoryCache
}
