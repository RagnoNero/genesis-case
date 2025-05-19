package dto

type FrequencyDto int

const (
	Invalid FrequencyDto = -1
	Hourly               = iota
	Daily
)
