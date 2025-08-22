package api

import (
	"context"
	"time"

	domain_exchange "exchange-rate-service/internal/domain/exchange"
)

type CompositeRepository struct {
	exchangeRateRepos []domain_exchange.ExchangeRateExternalRepository
}

func NewCompositeRepository(
	exchangeRateRepos ...domain_exchange.ExchangeRateExternalRepository,
) domain_exchange.ExchangeRateExternalRepository {
	return &CompositeRepository{
		exchangeRateRepos: exchangeRateRepos,
	}
}

func (c *CompositeRepository) GetLatestRate(ctx context.Context, fromCurrency string) (*domain_exchange.ExchangeRate, error) {
	repo := c.SelectRepo(fromCurrency)
	return repo.GetLatestRate(ctx, fromCurrency)
}

func (c *CompositeRepository) GetRateByDate(ctx context.Context, fromCurrency, toCurrency string, date time.Time) (*domain_exchange.ExchangeRate, error) {
	repo := c.SelectRepo(fromCurrency)
	return repo.GetRateByDate(ctx, fromCurrency, toCurrency, date)
}

func (c *CompositeRepository) GetRatesForDateRange(ctx context.Context, fromCurrency, toCurrency string, startDate, endDate time.Time) ([]*domain_exchange.ExchangeRate, error) {
	repo := c.SelectRepo(fromCurrency)
	return repo.GetRatesForDateRange(ctx, fromCurrency, toCurrency, startDate, endDate)
}

func (c *CompositeRepository) SelectRepo(fromCurrency string) domain_exchange.ExchangeRateExternalRepository {
	// 0 -> Fiat
	// 1 -> Crypto
	// 2 -> Mock

	if isCryptoSymbol(fromCurrency) {
		return c.exchangeRateRepos[1]
	}

	return c.exchangeRateRepos[0]
}

func isCryptoSymbol(symbol string) bool {
	return domain_exchange.SupportedCurrencies[symbol].Type == "crypto"
}
