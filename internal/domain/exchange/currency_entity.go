package domain_exchange

import (
	"errors"
	"fmt"
	"time"
)

type Currency struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

var SupportedCurrencies = map[string]Currency{
	"USD": {
		Code:   "USD",
		Name:   "United States Dollar",
		Symbol: "$",
	},
	"INR": {
		Code:   "INR",
		Name:   "Indian Rupee",
		Symbol: "₹",
	},
	"EUR": {
		Code:   "EUR",
		Name:   "Euro",
		Symbol: "€",
	},
	"JPY": {
		Code:   "JPY",
		Name:   "Japanese Yen",
		Symbol: "¥",
	},
	"GBP": {
		Code:   "GBP",
		Name:   "British Pound Sterling",
		Symbol: "£",
	},
}

type ExchangeRate struct {
	Result          string             `json:"result"`
	BaseCode        string             `json:"base_code"`
	ConversionRates map[string]float64 `json:"conversion_rates"`
	FetchedAt       time.Time          `json:"-"`
}

func (c *Currency) ValidateCurrencies(from, to string) error {
	if from == "" || to == "" {
		return errors.New("from and to currencies are required")
	}

	if _, exists := SupportedCurrencies[from]; !exists {
		return fmt.Errorf("currency %s is not supported", from)
	}

	if _, exists := SupportedCurrencies[to]; !exists {
		return fmt.Errorf("currency %s is not supported", to)
	}

	return nil
}
