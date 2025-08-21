package domain_exchange

import (
	"context"
	"time"
)

// TODO : Currently used by both infra and domain layer. Seperation needed

type ExchangeRateExternalRepository interface {
	GetLatestRate(ctx context.Context, fromCurrency string) (*ExchangeRate, error)
	GetRateByDate(ctx context.Context, fromCurrency, toCurrency string, date time.Time) (*ExchangeRate, error)
	GetRatesForDateRange(ctx context.Context, fromCurrency, toCurrency string, startDate, endDate time.Time) ([]*ExchangeRate, error)
}

type ExchangeRateCacheRepository interface {
	StoreRate(ctx context.Context, rate *ExchangeRate) error
	GetCachedRate(ctx context.Context, fromCurrency, toCurrency string, date time.Time) (*ExchangeRate, error)
	CacheRate(ctx context.Context, rate *ExchangeRate, ttl time.Duration) error
}
