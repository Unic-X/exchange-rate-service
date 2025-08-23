package exchange

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	domain_exchange "exchange-rate-service/internal/domain/exchange"
	"exchange-rate-service/pkg/logger"
)

func NewExchangeRateUsecase(usecase domain_exchange.ExchangeRateUsercase) domain_exchange.ExchangeRateUsercase {
	return usecase
}

type exchangeRateUseCase struct {
	externalRepo      domain_exchange.ExchangeRateExternalRepository
	cacheRepo         domain_exchange.ExchangeRateCacheRepository
	maxHistoricalDays int
}

func NewExchangeRateUseCase(
	externalRepo domain_exchange.ExchangeRateExternalRepository,
	cacheRepo domain_exchange.ExchangeRateCacheRepository,
	maxHistoricalDays int,
) domain_exchange.ExchangeRateUsercase {
	return &exchangeRateUseCase{
		externalRepo:      externalRepo,
		cacheRepo:         cacheRepo,
		maxHistoricalDays: maxHistoricalDays,
	}
}

func (s *exchangeRateUseCase) ValidateCurrencies(from, to string) error {
	if from == "" || to == "" {
		return errors.New("from and to currencies are required")
	}
	if _, exists := domain_exchange.SupportedCurrencies[from]; !exists {
		return fmt.Errorf("currency %s is not supported", from)
	}
	if _, exists := domain_exchange.SupportedCurrencies[to]; !exists {
		return fmt.Errorf("currency %s is not supported", to)
	}
	return nil
}

func (s *exchangeRateUseCase) ValidateDate(date time.Time, maxHistoricalDays int) error {
	now := time.Now()
	maxPastDate := now.AddDate(0, 0, -maxHistoricalDays)
	if date.After(now) {
		return errors.New("date cannot be in the future")
	}
	if date.Before(maxPastDate) {
		return fmt.Errorf("date cannot be older than %d days", maxHistoricalDays)
	}
	return nil
}

func (s *exchangeRateUseCase) GetLatestRate(ctx context.Context, from, to string) (float64, error) {
	if err := s.ValidateCurrencies(from, to); err != nil {
		return 0, err
	}
	if cachedRate, err := s.cacheRepo.GetCachedRate(ctx, from, to, time.Now()); err == nil && cachedRate != nil {
		if rate, exists := cachedRate.ConversionRates[to]; exists {
			logger.Infof("Cache hit for latest rate %s to %s", from, to)
			return rate, nil
		}
	}
	rate, err := s.externalRepo.GetLatestRate(ctx, from)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch latest rate: %w", err)
	}
	if err := s.cacheRepo.StoreRate(ctx, rate); err != nil {
		logger.Errorf("Failed to cache rate: %v", err)
	}
	if conversionRate, exists := rate.ConversionRates[to]; exists {
		return conversionRate, nil
	}
	return 0, fmt.Errorf("conversion rate from %s to %s not found", from, to)
}

func (s *exchangeRateUseCase) getHistoricalRate(ctx context.Context, from, to string, date time.Time) (float64, error) {
	if err := s.ValidateCurrencies(from, to); err != nil {
		return 0, err
	}
	if err := s.ValidateDate(date, s.maxHistoricalDays); err != nil {
		return 0, err
	}
	if cachedRate, err := s.cacheRepo.GetCachedRate(ctx, from, to, date); err == nil && cachedRate != nil {
		if rate, exists := cachedRate.ConversionRates[to]; exists {
			logger.Infof("Cache hit for historical rate %s to %s on %s", from, to, date.Format("2006-01-02"))
			return rate, nil
		}
	}
	rate, err := s.externalRepo.GetRateByDate(ctx, from, to, date)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch historical rate: %w", err)
	}
	if err := s.cacheRepo.StoreRate(ctx, rate); err != nil {
		logger.Errorf("Failed to cache historical rate: %v", err)
	}
	if conversionRate, exists := rate.ConversionRates[to]; exists {
		return conversionRate, nil
	}
	return 0, fmt.Errorf("conversion rate from %s to %s not found for date %s", from, to, date.Format("2006-01-02"))
}

func (s *exchangeRateUseCase) resolveRateByDate(ctx context.Context, from, to string, targetDate time.Time) (float64, error) {
	if targetDate.IsZero() {
		return s.GetLatestRate(ctx, from, to)
	}
	return s.getHistoricalRate(ctx, from, to, targetDate)
}

func (s *exchangeRateUseCase) ConvertAmount(ctx context.Context, from, to string, amount float64, fromDate, toDate time.Time) (float64, float64, float64, float64, error) {
	if amount <= 0 {
		return 0, 0, 0, 0, errors.New("amount must be greater than 0")
	}

	fromRate, err := s.resolveRateByDate(ctx, from, to, fromDate)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	toRate, err := s.resolveRateByDate(ctx, from, to, toDate)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	convertedAmountAtFromDate := amount * fromRate
	convertedAmountAtToDate := amount * toRate
	return convertedAmountAtFromDate, convertedAmountAtToDate, fromRate, toRate, nil
}

func (s *exchangeRateUseCase) RefreshRates(ctx context.Context) error {
	logger.Info("Starting rate refresh for all supported currencies")
	var wg sync.WaitGroup
	for baseCurrency := range domain_exchange.SupportedCurrencies {
		wg.Add(1)
		go func() {
			rate, err := s.externalRepo.GetLatestRate(ctx, baseCurrency)
			if err != nil {
				logger.Errorf("Failed to refresh rate for %s: %v", baseCurrency, err)
				wg.Done()
				return
			}
			if err := s.cacheRepo.StoreRate(ctx, rate); err != nil {
				logger.Errorf("Failed to cache refreshed rate for %s: %v", baseCurrency, err)
			}
			wg.Done()
		}()
	}
	logger.Info("Rate refresh completed")
	return nil
}
