package domain_exchange

import (
	"context"
	"time"
)

type ExchangeRateUsercase interface {
	ConvertAmount(ctx context.Context, from, to string, amount float64, fromDate, toDate time.Time) (float64, float64, float64, float64, error)
	GetLatestRate(ctx context.Context, from, to string) (float64, error)
	RefreshRates(ctx context.Context) error
	ValidateCurrencies(from, to string) error
	ValidateDate(date time.Time, maxHistoricalDays int) error
}
