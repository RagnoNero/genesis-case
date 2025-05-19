package dto

import "time"

type SubscriptionDto struct {
	Email        string
	City         string
	Frequency    FrequencyDto
	Confirmed    bool
	Token        string
	SubscribedAt time.Time
	LastSentAt   time.Time
}

func NewSubscriptionDto(email string, city string, freq FrequencyDto, token string) *SubscriptionDto {
	return &SubscriptionDto{
		Email:        email,
		City:         city,
		Frequency:    Daily,
		Confirmed:    false,
		Token:        token,
		SubscribedAt: time.Time{},
		LastSentAt:   time.Time{},
	}
}
