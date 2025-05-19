package sql

import "weather-subscription/sql/dto"

type ISubscriptionRepository interface {
	CreateSubscription(sub *dto.SubscriptionDto) error
	IsSubscribed(email string) bool
	ConfirmToken(token string) error
	UpdateLastSent(email string) error
	GetConfirmedSubscriptions() ([]dto.SubscriptionDto, error)
	GetByToken(token string) (*dto.SubscriptionDto, error)
	Unsubscribe(token string) error
}
