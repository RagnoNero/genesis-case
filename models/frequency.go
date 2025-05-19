package models

import (
	"fmt"
	"strings"
	"weather-subscription/sql/dto"
)

type Frequency int

const (
	Hourly Frequency = iota + 1
	Daily
)

func ParseFrequency(frequencyStr string) (Frequency, error) {
	switch strings.ToLower(frequencyStr) {
	case "hourly":
		return Hourly, nil
	case "daily":
		return Daily, nil
	default:
		return -1, fmt.Errorf("invalid frequency: %s", frequencyStr)
	}
}

func (f Frequency) ToDto() dto.FrequencyDto {
	switch f {
	case Hourly:
		return dto.Hourly
	case Daily:
		return dto.Daily
	default:
		return dto.Invalid
	}
}
