package inmemory

import (
	"context"
	"fmt"
	"time"

	"exchange-rate-service/internal/domain/entity"
	"exchange-rate-service/internal/domain/repository"
	"exchange-rate-service/internal/infra/cache"
	"exchange-rate-service/pkg/logger"
)

type inMemoryRepository struct {
	cache cache.Cache
}

func NewInMemoryRepository(cache cache.Cache) repository.ExchangeRateRepository {
	return &inMemoryRepository{
		cache: cache,
	}
}

func (r *inMemoryRepository) generateCacheKey(fromCurrency, toCurrency string, date time.Time) string {
	return fmt.Sprintf("rate:%s:%s:%s", fromCurrency, toCurrency, date.Format("2006-01-02"))
}

func (r *inMemoryRepository) GetLatestRate(ctx context.Context, fromCurrency, toCurrency string) (*entity.ExchangeRate, error) {
	key := r.generateCacheKey(fromCurrency, toCurrency, time.Now())
	
	if value, exists := r.cache.Get(key); exists {
		if rate, ok := value.(*entity.ExchangeRate); ok {
			logger.Infof("Cache hit for latest rate %s to %s", fromCurrency, toCurrency)
			return rate, nil
		}
	}
	
	return nil, fmt.Errorf("rate not found in cache")
}

func (r *inMemoryRepository) GetRateByDate(ctx context.Context, fromCurrency, toCurrency string, date time.Time) (*entity.ExchangeRate, error) {
	key := r.generateCacheKey(fromCurrency, toCurrency, date)
	
	if value, exists := r.cache.Get(key); exists {
		if rate, ok := value.(*entity.ExchangeRate); ok {
			logger.Infof("Cache hit for historical rate %s to %s on %s", fromCurrency, toCurrency, date.Format("2006-01-02"))
			return rate, nil
		}
	}
	
	return nil, fmt.Errorf("rate not found in cache for date %s", date.Format("2006-01-02"))
}

func (r *inMemoryRepository) GetRatesForDateRange(ctx context.Context, fromCurrency, toCurrency string, startDate, endDate time.Time) ([]*entity.ExchangeRate, error) {
	var rates []*entity.ExchangeRate
	
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		if rate, err := r.GetRateByDate(ctx, fromCurrency, toCurrency, d); err == nil {
			rates = append(rates, rate)
		}
	}
	
	return rates, nil
}

func (r *inMemoryRepository) StoreRate(ctx context.Context, rate *entity.ExchangeRate) error {
	if rate == nil {
		return fmt.Errorf("rate cannot be nil")
	}

	// Store rate for each supported currency pair
	for toCurrency := range rate.ConversionRates {
		key := r.generateCacheKey(rate.BaseCode, toCurrency, time.Now())
		if err := r.cache.Set(key, rate, 1*time.Hour); err != nil {
			return fmt.Errorf("failed to store rate in cache: %w", err)
		}
	}
	
	logger.Infof("Stored rate for %s in cache", rate.BaseCode)
	return nil
}

func (r *inMemoryRepository) GetCachedRate(ctx context.Context, fromCurrency, toCurrency string, date time.Time) (*entity.ExchangeRate, error) {
	return r.GetRateByDate(ctx, fromCurrency, toCurrency, date)
}

func (r *inMemoryRepository) CacheRate(ctx context.Context, rate *entity.ExchangeRate, ttl time.Duration) error {
	if rate == nil {
		return fmt.Errorf("rate cannot be nil")
	}

	for toCurrency := range rate.ConversionRates {
		key := r.generateCacheKey(rate.BaseCode, toCurrency, time.Now())
		if err := r.cache.Set(key, rate, ttl); err != nil {
			return fmt.Errorf("failed to cache rate: %w", err)
		}
	}
	
	return nil
}

func (r *inMemoryRepository) RefreshLatestRates(ctx context.Context) error {
	// In-memory repository doesn't need to refresh rates
	// This would be handled by the service layer
	return nil
}