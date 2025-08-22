package di

import (
	"context"
	"exchange-rate-service/internal/config"
	"exchange-rate-service/internal/delivery/http/handler"
	domain_exchange "exchange-rate-service/internal/domain/exchange"
	"exchange-rate-service/internal/infra/cache"
	"exchange-rate-service/internal/infra/http_client"
	"exchange-rate-service/internal/infra/repository/api"
	"exchange-rate-service/internal/infra/repository/inmemory"
	"exchange-rate-service/internal/infra/repository/mock"
	"exchange-rate-service/internal/usecase"
	"exchange-rate-service/pkg/logger"
	"time"
)

// Container holds all dependencies
type Container struct {
	Config                *config.Config
	HTTPClient            http_client.HTTPClient
	Cache                 cache.Cache
	ExternalAPIRepository domain_exchange.ExchangeRateExternalRepository
	InMemoryRepository    domain_exchange.ExchangeRateCacheRepository
	ExchangeRateUseCase   domain_exchange.ExchangeRateUsercase
	ExchangeRateHandler   *handler.ExchangeRateHandler
	MockRepository        domain_exchange.ExchangeRateExternalRepository
}

// NewContainer creates and wires all dependencies
func NewContainer(ctx context.Context, cfg *config.Config) *Container {
	container := &Container{
		Config: cfg,
	}

	// Infrastructure layer
	container.HTTPClient = http_client.NewHTTPClient(cfg.FiatExternalAPI.Timeout)
	container.Cache = cache.NewInMemoryCache(cfg.Cache.TTL)

	// Repository layer
	fiatRepo := api.NewExternalAPIRepository(
		container.HTTPClient,
		cfg.FiatExternalAPI.BaseURL,
		cfg.FiatExternalAPI.Secret,
	)
	cryptoRepo := api.NewCryptoAPIRepository(
		container.HTTPClient,
		cfg.CryptoExternalAPI.BaseURL,
		cfg.CryptoExternalAPI.Secret,
	)
	container.InMemoryRepository = inmemory.NewInMemoryRepository(container.Cache)

	mockRepository := mock.NewMockExchangeRateRepository()

	container.ExternalAPIRepository = api.NewCompositeRepository(
		fiatRepo, cryptoRepo, mockRepository,
	)
	container.ExchangeRateUseCase = usecase.NewExchangeRateService(
		container.ExternalAPIRepository,
		container.InMemoryRepository,
		cfg.Cache.MaxHistoricalDays,
	)
	container.ExchangeRateHandler = handler.NewExchangeRateHandler(container.ExchangeRateUseCase)
	go startRateRefreshTicker(container, ctx, cfg.Cache.RefreshInterval)
	return container
}

func startRateRefreshTicker(container *Container, ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	logger.Infof("Starting rate refresh ticker with interval: %v", interval)

	if err := container.ExchangeRateUseCase.RefreshRates(ctx); err != nil {
		logger.Errorf("Initial rate refresh failed: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			logger.Info("Rate refresh ticker stopped")
			return
		case <-ticker.C:
			logger.Info("Refreshing exchange rates...")
			if err := container.ExchangeRateUseCase.RefreshRates(ctx); err != nil {
				logger.Errorf("Rate refresh failed: %v", err)
			} else {
				logger.Info("Exchange rates refreshed successfully")
			}
		}
	}
}
