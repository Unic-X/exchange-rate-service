package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"time"

	domain_exchange "exchange-rate-service/internal/domain/exchange"
	"exchange-rate-service/internal/infra/http_client"
	"exchange-rate-service/pkg/logger"
)

type coinlayerLiveResp struct {
	Success   bool               `json:"success"`
	Timestamp int64              `json:"timestamp"`
	Target    string             `json:"target"`
	Rates     map[string]float64 `json:"rates"`
}

type cryptoAPIRepository struct {
	httpClient http_client.HTTPClient
	baseURL    string
	apiKey     string
}

func NewCryptoAPIRepository(httpClient http_client.HTTPClient, baseURL, apiKey string) domain_exchange.ExchangeRateExternalRepository {
	return &cryptoAPIRepository{
		httpClient: httpClient,
		baseURL:    baseURL,
		apiKey:     apiKey,
	}
}

func (r *cryptoAPIRepository) GetLatestRate(ctx context.Context, fromCurrency string) (*domain_exchange.ExchangeRate, error) {
	q := url.Values{}
	q.Set("access_key", r.apiKey)

	endpoint := fmt.Sprintf("%s/live?%s", r.baseURL, q.Encode())
	resp, err := r.httpClient.Get(ctx, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("coinlayer latest request failed: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("coinlayer status %d: %s", resp.StatusCode, string(body))
	}

	var res coinlayerLiveResp
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("coinlayer unmarshal: %w", err)
	}
	if !res.Success {
		return nil, fmt.Errorf("coinlayer returned success=false")
	}

	// Coinlayer returns a price table with target USD: Rates[symbol] = price in USD for 1 unit of symbol.
	// Normalize into conversion table (target-per-from) using USD as pivot:
	// rate(from -> x) = USDPrice(from) / USDPrice(x). For USD, USDPrice(USD) = 1, so rate(USD->x) = 1/USDPrice(x).
	usdPrices := res.Rates
	// Ensure USD is present and equals 1 if missing
	if _, ok := usdPrices["USD"]; !ok {
		usdPrices["USD"] = 1
	}

	fromUSDPrice, ok := usdPrices[fromCurrency]
	if !ok || fromUSDPrice == 0 {
		return nil, fmt.Errorf("coinlayer price in USD for %s not available", fromCurrency)
	}

	conversion := make(map[string]float64, len(usdPrices))
	for symbol, usdPrice := range usdPrices {
		if usdPrice == 0 {
			continue
		}
		conversion[symbol] = fromUSDPrice / usdPrice
	}
	// Identity
	conversion[fromCurrency] = 1

	rate := &domain_exchange.ExchangeRate{
		Result:          "success",
		BaseCode:        fromCurrency,
		ConversionRates: conversion,
		FetchedAt:       time.Now(),
	}
	logger.Infof("Fetched latest table via coinlayer; base=%s, target=%s (USD pivot)", fromCurrency, res.Target)
	return rate, nil
}

func (r *cryptoAPIRepository) GetRateByDate(ctx context.Context, fromCurrency, toCurrency string, date time.Time) (*domain_exchange.ExchangeRate, error) {
	dateStr := date.Format("2006-01-02")

	url := fmt.Sprintf("%s/%s?access_key=%s", r.baseURL, dateStr, r.apiKey)

	resp, err := r.httpClient.Get(ctx, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch historical rate: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var res coinlayerLiveResp
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("coinlayer unmarshal: %w", err)
	}
	if !res.Success {
		return nil, fmt.Errorf("coinlayer returned success=false")
	}

	return &domain_exchange.ExchangeRate{
		Result:          "success",
		BaseCode:        fromCurrency,
		ConversionRates: res.Rates,
		FetchedAt:       time.Now(),
	}, nil
}

func (r *cryptoAPIRepository) GetRatesForDateRange(ctx context.Context, fromCurrency, toCurrency string, startDate, endDate time.Time) ([]*domain_exchange.ExchangeRate, error) {
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
