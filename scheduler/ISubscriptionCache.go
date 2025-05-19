package scheduler

import (
	"time"
	"weather-subscription/sql/dto"
)

type ISubscriptionCache interface {
	ReloadCache() error
	AddToCache(token string)
	RemoveFromCache(token string)
	GetAll() []dto.SubscriptionDto
	UpdateLastSent(email string, lastSent time.Time)
}
