package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"exchange-rate-service/internal/domain/entity"
	"exchange-rate-service/internal/domain/repository"
	"exchange-rate-service/internal/infra/http_client"
	"exchange-rate-service/pkg/logger"
)

type externalAPIRepository struct {
	httpClient http_client.HTTPClient
	baseURL    string
	apiKey     string
}

func NewExternalAPIRepository(httpClient http_client.HTTPClient, baseURL, apiKey string) repository.ExchangeRateRepository {
	return &externalAPIRepository{
		httpClient: httpClient,
		baseURL:    baseURL,
		apiKey:     apiKey,
	}
}

func (r *externalAPIRepository) GetLatestRate(ctx context.Context, fromCurrency, toCurrency string) (*entity.ExchangeRate, error) {
	url := fmt.Sprintf("%s/%s/latest/%s", r.baseURL, r.apiKey, fromCurrency)

	resp, err := r.httpClient.Get(ctx, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest rate: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var rate entity.ExchangeRate
	if err := json.Unmarshal(body, &rate); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	rate.FetchedAt = time.Now()
	logger.Infof("Fetched latest rate for %s from external API", fromCurrency)

	return &rate, nil
}

func (r *externalAPIRepository) GetRateByDate(ctx context.Context, fromCurrency, toCurrency string, date time.Time) (*entity.ExchangeRate, error) {
	dateStr := date.Format("2006/01/02")

	// Below API requires paid version so it should work in theory.
	url := fmt.Sprintf("%s/%s/history/%s/%s", r.baseURL, r.apiKey, fromCurrency, dateStr)

	resp, err := r.httpClient.Get(ctx, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch historical rate: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var rate entity.ExchangeRate
	if err := json.Unmarshal(body, &rate); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	rate.FetchedAt = time.Now()
	logger.Infof("Fetched historical rate for %s on %s from external API", fromCurrency, dateStr)

	return &rate, nil
}

func (r *externalAPIRepository) GetRatesForDateRange(ctx context.Context, fromCurrency, toCurrency string, startDate, endDate time.Time) ([]*entity.ExchangeRate, error) {
	var rates []*entity.ExchangeRate

	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		rate, err := r.GetRateByDate(ctx, fromCurrency, toCurrency, d)
		if err != nil {
			logger.Errorf("Failed to fetch rate for %s: %v", d.Format("2006-01-02"), err)
			continue
		}
		rates = append(rates, rate)
	}

	return rates, nil
}

func (r *externalAPIRepository) StoreRate(ctx context.Context, rate *entity.ExchangeRate) error {
	return nil
}

func (r *externalAPIRepository) GetCachedRate(ctx context.Context, fromCurrency, toCurrency string, date time.Time) (*entity.ExchangeRate, error) {
	return nil, fmt.Errorf("no cached rate available")
}

func (r *externalAPIRepository) CacheRate(ctx context.Context, rate *entity.ExchangeRate, ttl time.Duration) error {
	return nil
}

func (r *externalAPIRepository) RefreshLatestRates(ctx context.Context) error {
	return nil
}
