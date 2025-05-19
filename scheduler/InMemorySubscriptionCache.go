package scheduler

import (
	"sync"
	"time"
	"weather-subscription/sql"
	"weather-subscription/sql/dto"
)

type InMemorySubscriptionCache struct {
	mu             sync.RWMutex
	cache          map[string]dto.SubscriptionDto
	repo           sql.ISubscriptionRepository
	lastReload     time.Time
	reloadInterval time.Duration
}

func NewInMemorySubscriptionCache(repo sql.ISubscriptionRepository) *InMemorySubscriptionCache {
	return &InMemorySubscriptionCache{
		cache:          make(map[string]dto.SubscriptionDto),
		repo:           repo,
		reloadInterval: 24 * time.Hour,
	}
}

func (c *InMemorySubscriptionCache) ReloadCache() error {
	subs, err := c.repo.GetConfirmedSubscriptions()
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[string]dto.SubscriptionDto)
	for _, sub := range subs {
		c.cache[sub.Email] = sub
	}

	c.lastReload = time.Now()
	return nil
}

func (c *InMemorySubscriptionCache) AddToCache(token string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	sub, _ := c.repo.GetByToken(token)
	c.cache[sub.Email] = *sub
}

func (c *InMemorySubscriptionCache) RemoveFromCache(token string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for email, sub := range c.cache {
		if sub.Token == token {
			delete(c.cache, email)
			break
		}
	}
}

func (c *InMemorySubscriptionCache) UpdateLastSent(email string, lastSent time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if sub, ok := c.cache[email]; ok {
		sub.LastSentAt = lastSent
		c.cache[email] = sub
	}
}

func (c *InMemorySubscriptionCache) GetAll() []dto.SubscriptionDto {
	c.mu.RLock()
	defer c.mu.RUnlock()
	var result []dto.SubscriptionDto
	for _, sub := range c.cache {
		result = append(result, sub)
	}
	return result
}
