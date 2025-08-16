package usecase

import (
	"context"
	"testing"
	"time"

	"exchange-rate-service/internal/domain/service"
)

// testService implements service.ExchangeRateService
type testService struct {
	latestRate float64
	histRate   float64
	convRate   float64
	convAmt    float64
	retErr     error
}

func (m *testService) GetLatestRate(ctx context.Context, from, to string) (float64, error) {
	return m.latestRate, m.retErr
}
func (m *testService) GetHistoricalRate(ctx context.Context, from, to string, date time.Time) (float64, error) {
	return m.histRate, m.retErr
}
func (m *testService) ConvertAmount(ctx context.Context, from, to string, amount float64, date time.Time) (float64, float64, error) {
	return m.convRate, m.convAmt, m.retErr
}
func (m *testService) RefreshRates(ctx context.Context) error                   { return m.retErr }
func (m *testService) ValidateCurrencies(from, to string) error                 { return m.retErr }
func (m *testService) ValidateDate(date time.Time, maxHistoricalDays int) error { return m.retErr }

var _ service.ExchangeRateService = (*testService)(nil)

func TestUsecase_DelegatesToService(t *testing.T) {
	ms := &testService{latestRate: 1.23, histRate: 0.99, convRate: 2, convAmt: 20}
	u := NewExchangeRateUsecase(ms)

	if r, err := u.GetLatestRate("EUR", "USD"); err != nil || r != 1.23 {
		t.Fatalf("GetLatestRate got r=%v err=%v", r, err)
	}

	date := time.Now().AddDate(0, 0, -1)
	if r, err := u.GetHistoricalRate("EUR", "USD", date); err != nil || r != 0.99 {
		t.Fatalf("GetHistoricalRate got r=%v err=%v", r, err)
	}

	if rate, amt, err := u.ConvertAmount("EUR", "USD", 10, time.Now()); err != nil || rate != 2 || amt != 20 {
		t.Fatalf("ConvertAmount got rate=%v amt=%v err=%v", rate, amt, err)
	}
}
