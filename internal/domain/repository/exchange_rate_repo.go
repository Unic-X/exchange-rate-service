package repository

import (
	"context"
	"exchange-rate-service/internal/domain/entity"
	"time"
)

// TODO : Currently used by both infra and domain layer. Seperation needed

type ExchangeRateRepository interface {
	GetLatestRate(ctx context.Context, fromCurrency, toCurrency string) (*entity.ExchangeRate, error)
	GetRateByDate(ctx context.Context, fromCurrency, toCurrency string, date time.Time) (*entity.ExchangeRate, error)
	GetRatesForDateRange(ctx context.Context, fromCurrency, toCurrency string, startDate, endDate time.Time) ([]*entity.ExchangeRate, error)

	// TODO : Remove from infra layer
	StoreRate(ctx context.Context, rate *entity.ExchangeRate) error
	GetCachedRate(ctx context.Context, fromCurrency, toCurrency string, date time.Time) (*entity.ExchangeRate, error)
	CacheRate(ctx context.Context, rate *entity.ExchangeRate, ttl time.Duration) error
	RefreshLatestRates(ctx context.Context) error
}
