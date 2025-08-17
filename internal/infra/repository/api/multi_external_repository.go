package api

import (
	"context"
	"time"

	"exchange-rate-service/internal/domain/entity"
	"exchange-rate-service/internal/domain/repository"
)

type multiExternalRepository struct {
	fiatRepo   repository.ExchangeRateRepository
	cryptoRepo repository.ExchangeRateRepository
}

func NewMultiExternalRepository(fiat repository.ExchangeRateRepository, crypto repository.ExchangeRateRepository) repository.ExchangeRateRepository {
	return &multiExternalRepository{fiatRepo: fiat, cryptoRepo: crypto}
}

func (m *multiExternalRepository) GetLatestRate(ctx context.Context, fromCurrency, toCurrency string) (*entity.ExchangeRate, error) {
	if isCryptoSymbol(fromCurrency) || isCryptoSymbol(toCurrency) {
		return m.cryptoRepo.GetLatestRate(ctx, fromCurrency, toCurrency)
	}
	return m.fiatRepo.GetLatestRate(ctx, fromCurrency, toCurrency)
}

func (m *multiExternalRepository) GetRateByDate(ctx context.Context, fromCurrency, toCurrency string, date time.Time) (*entity.ExchangeRate, error) {
	if isCryptoSymbol(fromCurrency) || isCryptoSymbol(toCurrency) {
		return m.cryptoRepo.GetRateByDate(ctx, fromCurrency, toCurrency, date)
	}
	return m.fiatRepo.GetRateByDate(ctx, fromCurrency, toCurrency, date)
}

func (m *multiExternalRepository) GetRatesForDateRange(ctx context.Context, fromCurrency, toCurrency string, startDate, endDate time.Time) ([]*entity.ExchangeRate, error) {
	if isCryptoSymbol(fromCurrency) || isCryptoSymbol(toCurrency) {
		return m.cryptoRepo.GetRatesForDateRange(ctx, fromCurrency, toCurrency, startDate, endDate)
	}
	return m.fiatRepo.GetRatesForDateRange(ctx, fromCurrency, toCurrency, startDate, endDate)
}

func (m *multiExternalRepository) StoreRate(ctx context.Context, rate *entity.ExchangeRate) error { return nil }
func (m *multiExternalRepository) GetCachedRate(ctx context.Context, fromCurrency, toCurrency string, date time.Time) (*entity.ExchangeRate, error) {
	return nil, nil
}
func (m *multiExternalRepository) CacheRate(ctx context.Context, rate *entity.ExchangeRate, ttl time.Duration) error { return nil }
func (m *multiExternalRepository) RefreshLatestRates(ctx context.Context) error { return nil }
