package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"exchange-rate-service/internal/domain/entity"
	"exchange-rate-service/internal/domain/repository"
	"exchange-rate-service/pkg/logger"
)

type ExchangeRateService interface {
	GetLatestRate(ctx context.Context, from, to string) (float64, error)
	GetHistoricalRate(ctx context.Context, from, to string, date time.Time) (float64, error)
	ConvertAmount(ctx context.Context, from, to string, amount float64, date time.Time) (float64, float64, error)
	RefreshRates(ctx context.Context) error
	ValidateCurrencies(from, to string) error
	ValidateDate(date time.Time, maxHistoricalDays int) error
}

type exchangeRateService struct {
	externalRepo      repository.ExchangeRateRepository
	cacheRepo         repository.ExchangeRateRepository
	maxHistoricalDays int
}

func NewExchangeRateService(
	externalRepo repository.ExchangeRateRepository,
	cacheRepo repository.ExchangeRateRepository,
	maxHistoricalDays int,
) ExchangeRateService {
	return &exchangeRateService{
		externalRepo:      externalRepo,
		cacheRepo:         cacheRepo,
		maxHistoricalDays: maxHistoricalDays,
	}
}

func (s *exchangeRateService) ValidateCurrencies(from, to string) error {
	if from == "" || to == "" {
		return errors.New("from and to currencies are required")
	}

	if _, exists := entity.SupportedCurrencies[from]; !exists {
		return fmt.Errorf("currency %s is not supported", from)
	}

	if _, exists := entity.SupportedCurrencies[to]; !exists {
		return fmt.Errorf("currency %s is not supported", to)
	}

	return nil
}

func (s *exchangeRateService) ValidateDate(date time.Time, maxHistoricalDays int) error {
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

func (s *exchangeRateService) GetLatestRate(ctx context.Context, from, to string) (float64, error) {
	if err := s.ValidateCurrencies(from, to); err != nil {
		return 0, err
	}

	// Try cache first
	if cachedRate, err := s.cacheRepo.GetLatestRate(ctx, from, to); err == nil && cachedRate != nil {
		if rate, exists := cachedRate.ConversionRates[to]; exists {
			logger.Infof("Cache hit for latest rate %s to %s", from, to)
			return rate, nil
		}
	}

	// Fetch from external API
	rate, err := s.externalRepo.GetLatestRate(ctx, from, to)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch latest rate: %w", err)
	}

	// Cache the result
	if err := s.cacheRepo.StoreRate(ctx, rate); err != nil {
		logger.Errorf("Failed to cache rate: %v", err)
	}

	if conversionRate, exists := rate.ConversionRates[to]; exists {
		return conversionRate, nil
	}

	return 0, fmt.Errorf("conversion rate from %s to %s not found", from, to)
}

func (s *exchangeRateService) GetHistoricalRate(ctx context.Context, from, to string, date time.Time) (float64, error) {
	if err := s.ValidateCurrencies(from, to); err != nil {
		return 0, err
	}

	if err := s.ValidateDate(date, s.maxHistoricalDays); err != nil {
		return 0, err
	}

	// Try cache first
	if cachedRate, err := s.cacheRepo.GetRateByDate(ctx, from, to, date); err == nil && cachedRate != nil {
		if rate, exists := cachedRate.ConversionRates[to]; exists {
			logger.Infof("Cache hit for historical rate %s to %s on %s", from, to, date.Format("2006-01-02"))
			return rate, nil
		}
	}

	// Fetch from external API
	rate, err := s.externalRepo.GetRateByDate(ctx, from, to, date)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch historical rate: %w", err)
	}

	// Cache the result
	if err := s.cacheRepo.StoreRate(ctx, rate); err != nil {
		logger.Errorf("Failed to cache historical rate: %v", err)
	}

	if conversionRate, exists := rate.ConversionRates[to]; exists {
		return conversionRate, nil
	}

	return 0, fmt.Errorf("conversion rate from %s to %s not found for date %s", from, to, date.Format("2006-01-02"))
}

func (s *exchangeRateService) ConvertAmount(ctx context.Context, from, to string, amount float64, date time.Time) (float64, float64, error) {
	if amount <= 0 {
		return 0, 0, errors.New("amount must be greater than 0")
	}

	var rate float64
	var err error

	// Use current date if date is today
	if date.Format("2006-01-02") == time.Now().Format("2006-01-02") {
		rate, err = s.GetLatestRate(ctx, from, to)
	} else {
		rate, err = s.GetHistoricalRate(ctx, from, to, date)
	}

	if err != nil {
		return 0, 0, err
	}

	convertedAmount := amount * rate
	return rate, convertedAmount, nil
}

func (s *exchangeRateService) RefreshRates(ctx context.Context) error {
	logger.Info("Starting rate refresh for all supported currencies")
	
	for baseCurrency := range entity.SupportedCurrencies {
		rate, err := s.externalRepo.GetLatestRate(ctx, baseCurrency, "USD")
		if err != nil {
			logger.Errorf("Failed to refresh rate for %s: %v", baseCurrency, err)
			continue
		}

		if err := s.cacheRepo.StoreRate(ctx, rate); err != nil {
			logger.Errorf("Failed to cache refreshed rate for %s: %v", baseCurrency, err)
		}
	}

	logger.Info("Rate refresh completed")
	return nil
}
