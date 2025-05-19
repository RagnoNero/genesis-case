package mocks

import (
	"time"
	"weather-subscription/sql/dto"
)

type DummyCache struct {
	subs map[string]dto.SubscriptionDto
}

func NewDummyCache() *DummyCache {
	return &DummyCache{
		subs: make(map[string]dto.SubscriptionDto),
	}
}

func (d *DummyCache) ReloadCache() error {
	return nil
}

func (d *DummyCache) AddToCache(token string) {
	d.subs[token] = dto.SubscriptionDto{
		Token: token,
	}
}

func (d *DummyCache) RemoveFromCache(token string) {
	delete(d.subs, token)
}

func (d *DummyCache) GetAll() []dto.SubscriptionDto {
	values := make([]dto.SubscriptionDto, 0, len(d.subs))
	for _, v := range d.subs {
		values = append(values, v)
	}
	return values
}

func (d *DummyCache) UpdateLastSent(email string, lastSent time.Time) {
	for token, sub := range d.subs {
		if sub.Email == email {
			sub.LastSentAt = lastSent
			d.subs[token] = sub
		}
	}
}
