package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"exchange-rate-service/internal/domain/entity"
)

// mockRepo implements repository.ExchangeRateRepository for tests
type mockRepo struct {
	latestResp *entity.ExchangeRate
	latestErr  error

	histResp *entity.ExchangeRate
	histErr  error

	stored []*entity.ExchangeRate
}

func (m *mockRepo) GetLatestRate(ctx context.Context, fromCurrency, toCurrency string) (*entity.ExchangeRate, error) {
	return m.latestResp, m.latestErr
}
func (m *mockRepo) GetRateByDate(ctx context.Context, fromCurrency, toCurrency string, date time.Time) (*entity.ExchangeRate, error) {
	return m.histResp, m.histErr
}
func (m *mockRepo) GetRatesForDateRange(ctx context.Context, fromCurrency, toCurrency string, startDate, endDate time.Time) ([]*entity.ExchangeRate, error) {
	return nil, nil
}
func (m *mockRepo) StoreRate(ctx context.Context, rate *entity.ExchangeRate) error {
	m.stored = append(m.stored, rate)
	return nil
}
func (m *mockRepo) GetCachedRate(ctx context.Context, fromCurrency, toCurrency string, date time.Time) (*entity.ExchangeRate, error) {
	return nil, nil
}
func (m *mockRepo) CacheRate(ctx context.Context, rate *entity.ExchangeRate, ttl time.Duration) error {
	return nil
}
func (m *mockRepo) RefreshLatestRates(ctx context.Context) error { return nil }

func mkRate(base string, rates map[string]float64) *entity.ExchangeRate {
	return &entity.ExchangeRate{BaseCode: base, ConversionRates: rates, FetchedAt: time.Now()}
}

func TestValidateCurrencies_Invalid(t *testing.T) {
	svc := NewExchangeRateService(&mockRepo{}, &mockRepo{}, 90)
	if err := svc.ValidateCurrencies("", "USD"); err == nil {
		t.Fatalf("expected error for empty from")
	}
	if err := svc.ValidateCurrencies("XXX", "USD"); err == nil {
		t.Fatalf("expected error for unsupported from currency")
	}
}

func TestValidateDate_FutureAndTooOld(t *testing.T) {
	svc := NewExchangeRateService(&mockRepo{}, &mockRepo{}, 10)
	if err := svc.ValidateDate(time.Now().Add(24*time.Hour), 10); err == nil {
		t.Fatalf("expected error for future date")
	}
	old := time.Now().AddDate(0, 0, -11)
	if err := svc.ValidateDate(old, 10); err == nil {
		t.Fatalf("expected error for too old date")
	}
}

func TestGetLatestRate_CacheHit(t *testing.T) {
	cache := &mockRepo{latestResp: mkRate("EUR", map[string]float64{"USD": 1.1})}
	ext := &mockRepo{}
	svc := NewExchangeRateService(ext, cache, 90)
	rate, err := svc.GetLatestRate(context.Background(), "EUR", "USD")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if rate != 1.1 {
		t.Fatalf("want 1.1 got %v", rate)
	}
}

func TestGetLatestRate_FallbackToExternalAndCache(t *testing.T) {
	cache := &mockRepo{latestResp: nil, latestErr: errors.New("miss")}
	ext := &mockRepo{latestResp: mkRate("EUR", map[string]float64{"USD": 1.2})}
	svc := NewExchangeRateService(ext, cache, 90)
	rate, err := svc.GetLatestRate(context.Background(), "EUR", "USD")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if rate != 1.2 {
		t.Fatalf("want 1.2 got %v", rate)
	}
	if len(cache.stored) != 1 {
		t.Fatalf("expected cached store, got %d", len(cache.stored))
	}
}

func TestGetHistoricalRate_ValidatesDateAndCacheHit(t *testing.T) {
	date := time.Now().AddDate(0, 0, -5)
	cache := &mockRepo{histResp: mkRate("EUR", map[string]float64{"USD": 1.05})}
	ext := &mockRepo{}
	svc := NewExchangeRateService(ext, cache, 90)
	rate, err := svc.GetHistoricalRate(context.Background(), "EUR", "USD", date)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if rate != 1.05 {
		t.Fatalf("want 1.05 got %v", rate)
	}
}

func TestConvertAmount_TodayUsesLatest(t *testing.T) {
	cache := &mockRepo{latestResp: mkRate("EUR", map[string]float64{"USD": 2})}
	ext := &mockRepo{}
	svc := NewExchangeRateService(ext, cache, 90)
	rate, converted, err := svc.ConvertAmount(context.Background(), "EUR", "USD", 10, time.Now())
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if rate != 2 || converted != 20 {
		t.Fatalf("want rate=2 converted=20 got rate=%v converted=%v", rate, converted)
	}
}

func TestConvertAmount_InvalidAmount(t *testing.T) {
	svc := NewExchangeRateService(&mockRepo{}, &mockRepo{}, 90)
	if _, _, err := svc.ConvertAmount(context.Background(), "EUR", "USD", 0, time.Now()); err == nil {
		t.Fatalf("expected error for non-positive amount")
	}
}
