package mock

import (
	"context"
	"fmt"
	"time"

	"exchange-rate-service/internal/domain/entity"
	"exchange-rate-service/internal/domain/repository"
)

type MockExchangeRateRepository struct{}

func NewMockExchangeRateRepository() repository.ExchangeRateRepository {
	return &MockExchangeRateRepository{}
}

func (m *MockExchangeRateRepository) GetLatestRate(ctx context.Context, fromCurrency, toCurrency string) (*entity.ExchangeRate, error) {
	return &entity.ExchangeRate{
		Result:   "success",
		BaseCode: fromCurrency,
		ConversionRates: map[string]float64{
			toCurrency: 1.0,
		},
		FetchedAt: time.Now(),
	}, nil
}

func (m *MockExchangeRateRepository) GetRateByDate(ctx context.Context, fromCurrency, toCurrency string, date time.Time) (*entity.ExchangeRate, error) {
	d := date.UTC()
	key := d.Year()*10000 + int(d.Month())*100 + d.Day()
	variance := float64((key%7)-3) * 0.005
	base := 0.95
	rate := base + variance
	if rate <= 0 {
		rate = 0.9
	}

	return &entity.ExchangeRate{
		Result:   "success",
		BaseCode: fromCurrency,
		ConversionRates: map[string]float64{
			toCurrency: rate,
		},
		FetchedAt: time.Now(),
	}, nil
}

func (m *MockExchangeRateRepository) GetRatesForDateRange(ctx context.Context, fromCurrency, toCurrency string, startDate, endDate time.Time) ([]*entity.ExchangeRate, error) {
	var rates []*entity.ExchangeRate
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		r, _ := m.GetRateByDate(ctx, fromCurrency, toCurrency, d)
		rates = append(rates, r)
	}
	return rates, nil
}

func (m *MockExchangeRateRepository) StoreRate(ctx context.Context, rate *entity.ExchangeRate) error {
	return nil
}

func (m *MockExchangeRateRepository) GetCachedRate(ctx context.Context, fromCurrency, toCurrency string, date time.Time) (*entity.ExchangeRate, error) {
	return nil, fmt.Errorf("no cached rate available in mock")
}

func (m *MockExchangeRateRepository) CacheRate(ctx context.Context, rate *entity.ExchangeRate, ttl time.Duration) error {
	return nil
}

func (m *MockExchangeRateRepository) RefreshLatestRates(ctx context.Context) error {
	return nil
}
