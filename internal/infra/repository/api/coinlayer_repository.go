package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"time"

	"exchange-rate-service/internal/domain/entity"
	"exchange-rate-service/internal/domain/repository"
	"exchange-rate-service/internal/infra/http_client"
	"exchange-rate-service/pkg/logger"
)

type coinlayerRepository struct {
	httpClient http_client.HTTPClient
	baseURL    string
	apiKey     string
}

func NewCoinlayerRepository(httpClient http_client.HTTPClient, baseURL, apiKey string) repository.ExchangeRateRepository {
	return &coinlayerRepository{
		httpClient: httpClient,
		baseURL:    baseURL,
		apiKey:     apiKey,
	}
}

type coinlayerLiveResp struct {
	Success   bool               `json:"success"`
	Target    string             `json:"target"`
	Timestamp int64              `json:"timestamp"`
	Rates     map[string]float64 `json:"rates"`
}

func (r *coinlayerRepository) GetLatestRate(ctx context.Context, fromCurrency, toCurrency string) (*entity.ExchangeRate, error) {
	// Determine query shape
	target := ""
	symbols := []string{}

	fromIsCrypto := isCryptoSymbol(fromCurrency)
	toIsCrypto := isCryptoSymbol(toCurrency)

	if fromIsCrypto && !toIsCrypto {
		// Crypto -> Fiat: target = fiat (to), symbols = from
		target = toCurrency
		symbols = []string{fromCurrency}
	} else if !fromIsCrypto && toIsCrypto {
		// Fiat -> Crypto: target = fiat (from), symbols = to
		target = fromCurrency
		symbols = []string{toCurrency}
	} else if fromIsCrypto && toIsCrypto {
		// Crypto -> Crypto: use USD pivot
		target = "USD"
		symbols = []string{fromCurrency, toCurrency}
	} else {
		return nil, fmt.Errorf("coinlayer supports crypto pairs; got fiat->fiat %s->%s", fromCurrency, toCurrency)
	}

	q := url.Values{}
	q.Set("access_key", r.apiKey)
	q.Set("target", target)
	q.Set("symbols", stringsJoin(symbols, ","))

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
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("coinlayer read body: %w", err)
	}
	var res coinlayerLiveResp
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("coinlayer unmarshal: %w", err)
	}
	if !res.Success {
		return nil, fmt.Errorf("coinlayer returned success=false")
	}

	conversion := make(map[string]float64)
	if fromIsCrypto && !toIsCrypto {
		price := res.Rates[fromCurrency]
		conversion[toCurrency] = price // 1 from equals price in to
	} else if !fromIsCrypto && toIsCrypto {
		price := res.Rates[toCurrency]
		if price == 0 {
			return nil, fmt.Errorf("coinlayer zero price for %s", toCurrency)
		}
		conversion[toCurrency] = 1.0 / price // 1 from fiat equals 1/price units of crypto
	} else {
		// crypto -> crypto via USD pivot
		priceFrom := res.Rates[fromCurrency]
		priceTo := res.Rates[toCurrency]
		if priceTo == 0 {
			return nil, fmt.Errorf("coinlayer zero price for %s", toCurrency)
		}
		conversion[toCurrency] = priceFrom / priceTo
	}

	rate := &entity.ExchangeRate{
		Result:          "success",
		BaseCode:        fromCurrency,
		ConversionRates: conversion,
		FetchedAt:       time.Now(),
	}
	logger.Infof("Fetched latest crypto rate %s->%s via coinlayer", fromCurrency, toCurrency)
	return rate, nil
}

func (r *coinlayerRepository) GetRateByDate(ctx context.Context, fromCurrency, toCurrency string, date time.Time) (*entity.ExchangeRate, error) {
	target := ""
	symbols := []string{}

	fromIsCrypto := isCryptoSymbol(fromCurrency)
	toIsCrypto := isCryptoSymbol(toCurrency)

	if fromIsCrypto && !toIsCrypto {
		target = toCurrency
		symbols = []string{fromCurrency}
	} else if !fromIsCrypto && toIsCrypto {
		target = fromCurrency
		symbols = []string{toCurrency}
	} else if fromIsCrypto && toIsCrypto {
		target = "USD"
		symbols = []string{fromCurrency, toCurrency}
	} else {
		return nil, fmt.Errorf("coinlayer supports crypto pairs; got fiat->fiat %s->%s", fromCurrency, toCurrency)
	}

	q := url.Values{}
	q.Set("access_key", r.apiKey)
	q.Set("target", target)
	q.Set("symbols", stringsJoin(symbols, ","))

	dateStr := date.Format("2006-01-02")
	endpoint := fmt.Sprintf("%s/%s?%s", r.baseURL, dateStr, q.Encode())
	resp, err := r.httpClient.Get(ctx, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("coinlayer historical request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("coinlayer status %d: %s", resp.StatusCode, string(body))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("coinlayer read body: %w", err)
	}
	var res coinlayerLiveResp
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("coinlayer unmarshal: %w", err)
	}
	if !res.Success {
		return nil, fmt.Errorf("coinlayer returned success=false")
	}

	conversion := make(map[string]float64)
	if fromIsCrypto && !toIsCrypto {
		price := res.Rates[fromCurrency]
		conversion[toCurrency] = price
	} else if !fromIsCrypto && toIsCrypto {
		price := res.Rates[toCurrency]
		if price == 0 {
			return nil, fmt.Errorf("coinlayer zero price for %s", toCurrency)
		}
		conversion[toCurrency] = 1.0 / price
	} else {
		priceFrom := res.Rates[fromCurrency]
		priceTo := res.Rates[toCurrency]
		if priceTo == 0 {
			return nil, fmt.Errorf("coinlayer zero price for %s", toCurrency)
		}
		conversion[toCurrency] = priceFrom / priceTo
	}

	rate := &entity.ExchangeRate{
		Result:          "success",
		BaseCode:        fromCurrency,
		ConversionRates: conversion,
		FetchedAt:       time.Now(),
	}
	logger.Infof("Fetched historical crypto rate %s->%s on %s via coinlayer", fromCurrency, toCurrency, dateStr)
	return rate, nil
}

func (r *coinlayerRepository) GetRatesForDateRange(ctx context.Context, fromCurrency, toCurrency string, startDate, endDate time.Time) ([]*entity.ExchangeRate, error) {
	var rates []*entity.ExchangeRate
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		rate, err := r.GetRateByDate(ctx, fromCurrency, toCurrency, d)
		if err != nil {
			logger.Errorf("coinlayer range fetch failed for %s: %v", d.Format("2006-01-02"), err)
			continue
		}
		rates = append(rates, rate)
	}
	return rates, nil
}

func (r *coinlayerRepository) StoreRate(ctx context.Context, rate *entity.ExchangeRate) error {
	return nil
}
func (r *coinlayerRepository) GetCachedRate(ctx context.Context, fromCurrency, toCurrency string, date time.Time) (*entity.ExchangeRate, error) {
	return nil, fmt.Errorf("no cached rate available")
}
func (r *coinlayerRepository) CacheRate(ctx context.Context, rate *entity.ExchangeRate, ttl time.Duration) error {
	return nil
}
func (r *coinlayerRepository) RefreshLatestRates(ctx context.Context) error { return nil }

func isCryptoSymbol(sym string) bool {
	return entity.SupportedCurrencies[sym].Type == "crypto"
}

func stringsJoin(arr []string, sep string) string {
	if len(arr) == 0 {
		return ""
	}
	out := arr[0]
	for i := 1; i < len(arr); i++ {
		out += sep + arr[i]
	}
	return out
}
