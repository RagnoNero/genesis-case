package mocks

import (
	"fmt"
	"time"
	"weather-subscription/sql/dto"
)

type DummySubRepo struct {
	Subs []dto.SubscriptionDto
}

func NewDummySubRepo() *DummySubRepo {
	return &DummySubRepo{}
}

func (r *DummySubRepo) CreateSubscription(sub *dto.SubscriptionDto) error {
	for _, s := range r.Subs {
		if s.Email == sub.Email {
			return fmt.Errorf("subscription already exists")
		}
	}
	r.Subs = append(r.Subs, *sub)
	return nil
}

func (r *DummySubRepo) IsSubscribed(email string) bool {
	for _, s := range r.Subs {
		if s.Email == email {
			return true
		}
	}
	return false
}

func (r *DummySubRepo) ConfirmToken(token string) error {
	for i, s := range r.Subs {
		if s.Token == token {
			r.Subs[i].Confirmed = true
			return nil
		}
	}
	return fmt.Errorf("token not found")
}

func (r *DummySubRepo) GetByToken(token string) (*dto.SubscriptionDto, error) {
	for _, sub := range r.Subs {
		if sub.Token == token {
			return &sub, nil
		}
	}
	return nil, nil
}

func (r *DummySubRepo) UpdateLastSent(email string) error {
	for i, s := range r.Subs {
		if s.Email == email {
			r.Subs[i].LastSentAt = time.Now()
			return nil
		}
	}
	return fmt.Errorf("email not found")
}

func (r *DummySubRepo) Unsubscribe(token string) error {
	for i, s := range r.Subs {
		if s.Token == token {
			r.Subs = append(r.Subs[:i], r.Subs[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("token not found")
}

func (r *DummySubRepo) GetConfirmedSubscriptions() ([]dto.SubscriptionDto, error) {
	return r.Subs, nil
}
