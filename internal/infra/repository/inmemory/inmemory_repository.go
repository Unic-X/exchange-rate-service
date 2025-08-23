package inmemory

import (
	"context"
	"fmt"
	"time"

	exchange "exchange-rate-service/internal/domain/exchange"
	"exchange-rate-service/pkg/logger"

	"exchange-rate-service/pkg/cache"
)

type inMemoryRepository struct {
	cache cache.Cache
}

func NewInMemoryRepository(cache cache.Cache) exchange.ExchangeRateCacheRepository {
	return &inMemoryRepository{
		cache: cache,
	}
}

func (r *inMemoryRepository) generateCacheKey(fromCurrency string, date time.Time) string {
	return fmt.Sprintf("rate:%s:%s", fromCurrency, date.Format("2006-01-02"))
}

func (r *inMemoryRepository) GetLatestRate(ctx context.Context, fromCurrency string) (*exchange.ExchangeRate, error) {
	key := r.generateCacheKey(fromCurrency, time.Now())

	if value, exists := r.cache.Get(key); exists {
		if rate, ok := value.(*exchange.ExchangeRate); ok {
			logger.Infof("Cache hit for latest base table %s", fromCurrency)
			return rate, nil
		}
	}

	return nil, fmt.Errorf("rate not found in cache")
}

func (r *inMemoryRepository) GetRateByDate(ctx context.Context, fromCurrency, toCurrency string, date time.Time) (*exchange.ExchangeRate, error) {
	key := r.generateCacheKey(fromCurrency, date)

	if value, exists := r.cache.Get(key); exists {
		if rate, ok := value.(*exchange.ExchangeRate); ok {
			logger.Infof("Cache hit for historical base table %s on %s", fromCurrency, date.Format("2006-01-02"))
			return rate, nil
		}
	}

	return nil, fmt.Errorf("rate not found in cache for date %s", date.Format("2006-01-02"))
}

func (r *inMemoryRepository) GetRatesForDateRange(ctx context.Context, fromCurrency, toCurrency string, startDate, endDate time.Time) ([]*exchange.ExchangeRate, error) {
	var rates []*exchange.ExchangeRate

	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		if rate, err := r.GetRateByDate(ctx, fromCurrency, toCurrency, d); err == nil {
			rates = append(rates, rate)
		}
	}

	return rates, nil
}

func (r *inMemoryRepository) StoreRate(ctx context.Context, rate *exchange.ExchangeRate) error {
	if rate == nil {
		return fmt.Errorf("rate cannot be nil")
	}

	key := r.generateCacheKey(rate.BaseCode, time.Now())
	if err := r.cache.Set(key, rate, 1*time.Hour); err != nil {
		return fmt.Errorf("failed to store rate in cache: %w", err)
	}
	logger.Infof("Stored rate for %s in cache", rate.BaseCode)
	return nil
}

func (r *inMemoryRepository) GetCachedRate(ctx context.Context, fromCurrency, toCurrency string, date time.Time) (*exchange.ExchangeRate, error) {
	return r.GetRateByDate(ctx, fromCurrency, toCurrency, date)
}

func (r *inMemoryRepository) CacheRate(ctx context.Context, rate *exchange.ExchangeRate, ttl time.Duration) error {
	if rate == nil {
		return fmt.Errorf("rate cannot be nil")
	}
	key := r.generateCacheKey(rate.BaseCode, time.Now())
	if err := r.cache.Set(key, rate, ttl); err != nil {
		return fmt.Errorf("failed to cache rate: %w", err)
	}

	return nil
}
