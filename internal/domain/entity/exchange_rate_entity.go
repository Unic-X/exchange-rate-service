package entity

import "time"

type ExchangeRate struct {
	Result          string             `json:"result"`
	BaseCode        string             `json:"base_code"`
	ConversionRates map[string]float64 `json:"conversion_rates"`
	FetchedAt       time.Time          `json:"-"`
}
