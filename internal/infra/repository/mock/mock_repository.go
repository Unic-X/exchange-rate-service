package mock

import (
	"context"
	"math/rand"
	"time"

	entity "exchange-rate-service/internal/domain/exchange"
)

type MockExchangeRateRepository struct {
	//Doesn't need anything``
}

func NewMockExchangeRateRepository() entity.ExchangeRateExternalRepository {
	return &MockExchangeRateRepository{}
}

func (m *MockExchangeRateRepository) GetLatestRate(ctx context.Context, fromCurrency string) (*entity.ExchangeRate, error) {
	conversionRates := make(map[string]float64)
	conversionRates[fromCurrency] = 1.0

	for currency := range entity.SupportedCurrencies {
		conversionRates[currency] = (rand.Float64() * 5)
	}

	return &entity.ExchangeRate{
		Result:          "success",
		BaseCode:        fromCurrency,
		ConversionRates: conversionRates,
		FetchedAt:       time.Now(),
	}, nil

	//HardCoded For now
}

func (m *MockExchangeRateRepository) GetRateByDate(ctx context.Context, fromCurrency, toCurrency string, date time.Time) (*entity.ExchangeRate, error) {
	d := date.UTC()
	key := d.Year()*10000 + int(d.Month())*100 + d.Day()*10
	variance := float64((key%7)-3) * 0.05
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
