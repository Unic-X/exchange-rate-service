package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	domain_exchange "exchange-rate-service/internal/domain/exchange"
	"exchange-rate-service/internal/infra/http_client"
	"exchange-rate-service/pkg/logger"
)

type externalAPIRepository struct {
	httpClient http_client.HTTPClient
	baseURL    string
	apiKey     string
}

func NewExternalAPIRepository(httpClient http_client.HTTPClient, baseURL, apiKey string) domain_exchange.ExchangeRateExternalRepository {
	return &externalAPIRepository{
		httpClient: httpClient,
		baseURL:    baseURL,
		apiKey:     apiKey,
	}
}

func (r *externalAPIRepository) GetLatestRate(ctx context.Context, fromCurrency string) (*domain_exchange.ExchangeRate, error) {
	url := fmt.Sprintf("%s/%s/latest/%s", r.baseURL, r.apiKey, fromCurrency)

	resp, err := r.httpClient.Get(ctx, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest rate: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}
	var rate domain_exchange.ExchangeRate
	if err := json.NewDecoder(resp.Body).Decode(&rate); err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	now := time.Now()
	rate.FetchedAt = now
	logger.Infof("Fetched latest rate for %s from external API", fromCurrency)

	return &rate, nil
}

func (r *externalAPIRepository) GetRateByDate(ctx context.Context, fromCurrency, toCurrency string, date time.Time) (*domain_exchange.ExchangeRate, error) {
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

	var rate domain_exchange.ExchangeRate
	if err := json.Unmarshal(body, &rate); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	rate.FetchedAt = time.Now()
	logger.Infof("Fetched historical rate for %s on %s from external API", fromCurrency, dateStr)

	return &rate, nil
}

func (r *externalAPIRepository) GetRatesForDateRange(ctx context.Context, fromCurrency, toCurrency string, startDate, endDate time.Time) ([]*domain_exchange.ExchangeRate, error) {
	var rates []*domain_exchange.ExchangeRate

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
