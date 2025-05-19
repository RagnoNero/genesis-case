package unit

import (
	"testing"
	"time"

	"weather-subscription/scheduler"
	"weather-subscription/sql/dto"
	"weather-subscription/tests/mocks"
)

func createTestSubscription(email string, confirmed bool) dto.SubscriptionDto {
	return dto.SubscriptionDto{
		Email:        email,
		City:         "Kyiv",
		Frequency:    1,
		Confirmed:    confirmed,
		Token:        "token-" + email,
		SubscribedAt: time.Now().Add(-time.Hour),
		LastSentAt:   time.Now().Add(-time.Minute),
	}
}

func TestReloadCache_LoadsSubscriptions(t *testing.T) {
	sub := createTestSubscription("user1@example.com", true)
	repo := &mocks.DummySubRepo{Subs: []dto.SubscriptionDto{sub}}

	cache := scheduler.NewInMemorySubscriptionCache(repo)
	err := cache.ReloadCache()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	all := cache.GetAll()
	if len(all) != 1 {
		t.Errorf("expected 1 subscription, got %d", len(all))
	}
	if all[0].Email != "user1@example.com" {
		t.Errorf("unexpected email: %s", all[0].Email)
	}
}

func TestAddToCache_AddsSubscription(t *testing.T) {
	sub := createTestSubscription("user2@example.com", true)
	repo := &mocks.DummySubRepo{Subs: []dto.SubscriptionDto{sub}}
	cache := scheduler.NewInMemorySubscriptionCache(repo)
	cache.AddToCache(sub.Token)

	all := cache.GetAll()
	if len(all) != 1 {
		t.Errorf("expected 1 subscription, got %d", len(all))
	}
}

func TestRemoveFromCache_RemovesSubscription(t *testing.T) {
	sub := createTestSubscription("user3@example.com", true)
	repo := &mocks.DummySubRepo{Subs: []dto.SubscriptionDto{sub}}

	cache := scheduler.NewInMemorySubscriptionCache(repo)
	cache.ReloadCache()
	cache.RemoveFromCache("token-user3@example.com")

	all := cache.GetAll()
	if len(all) != 0 {
		t.Errorf("expected cache to be empty, got %d", len(all))
	}
}

func TestUpdateLastSent_UpdatesTimestamp(t *testing.T) {
	sub := createTestSubscription("user4@example.com", true)
	repo := &mocks.DummySubRepo{Subs: []dto.SubscriptionDto{sub}}

	cache := scheduler.NewInMemorySubscriptionCache(repo)
	cache.ReloadCache()

	newTime := time.Now().Add(1 * time.Hour)
	cache.UpdateLastSent("user4@example.com", newTime)

	found := false
	for _, s := range cache.GetAll() {
		if s.Email == "user4@example.com" {
			found = true
			if !s.LastSentAt.Equal(newTime) {
				t.Errorf("expected last_sent_at to be %v, got %v", newTime, s.LastSentAt)
			}
		}
	}
	if !found {
		t.Error("subscription not found after update")
	}
}
