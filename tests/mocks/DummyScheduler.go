package mocks

import (
	"weather-subscription/scheduler"
)

type DummyScheduler struct {
	Cache scheduler.ISubscriptionCache
}

func NewDummyScheduler() scheduler.IScheduler {
	return &DummyScheduler{
		Cache: NewDummyCache(),
	}
}

func (s *DummyScheduler) Start() {

}

func (s *DummyScheduler) GetCache() scheduler.ISubscriptionCache {
	return s.Cache
}
