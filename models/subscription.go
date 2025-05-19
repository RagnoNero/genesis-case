package models

import (
	"weather-subscription/sql/dto"
)

type Subscription struct {
	Email     string    `json:"email"`
	City      string    `json:"city"`
	Frequency Frequency `json:"frequency"`
}

func (s *Subscription) ToDto() *dto.SubscriptionDto {
	return &dto.SubscriptionDto{
		Email:     s.Email,
		City:      s.City,
		Frequency: s.Frequency.ToDto(),
		Token:     "",
		Confirmed: false,
	}
}
