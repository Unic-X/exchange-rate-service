package di

import (
	"exchange-rate-service/internal/config"
	"exchange-rate-service/internal/delivery/http/handler"
	"exchange-rate-service/internal/domain/repository"
	"exchange-rate-service/internal/domain/service"
	"exchange-rate-service/internal/infra/cache"
	"exchange-rate-service/internal/infra/http_client"
	"exchange-rate-service/internal/infra/repository/api"
	"exchange-rate-service/internal/infra/repository/inmemory"
	"exchange-rate-service/internal/usecase"
)

// Container holds all dependencies
type Container struct {
	Config                *config.Config
	HTTPClient            http_client.HTTPClient
	Cache                 cache.Cache
	ExternalAPIRepository repository.ExchangeRateRepository
	InMemoryRepository    repository.ExchangeRateRepository
	ExchangeRateService   service.ExchangeRateService
	ExchangeRateUsecase   usecase.ExchangeRateUsecase
	ExchangeRateHandler   *handler.ExchangeRateHandler
}

// NewContainer creates and wires all dependencies
func NewContainer(cfg *config.Config) *Container {
	container := &Container{
		Config: cfg,
	}

	// Infrastructure layer
	container.HTTPClient = http_client.NewHTTPClient(cfg.FiatExternalAPI.Timeout)
	container.Cache = cache.NewInMemoryCache(cfg.Cache.TTL)

	fiatRepo := api.NewExternalAPIRepository(
		container.HTTPClient,
		cfg.FiatExternalAPI.BaseURL,
		cfg.FiatExternalAPI.Secret,
	)
	cryptoRepo := api.NewCoinlayerRepository(
		container.HTTPClient,
		cfg.CryptoExternalAPI.BaseURL,
		cfg.CryptoExternalAPI.Secret,
	)

	container.ExternalAPIRepository = api.NewMultiExternalRepository(fiatRepo, cryptoRepo)
	container.InMemoryRepository = inmemory.NewInMemoryRepository(container.Cache)

	// Service layer
	container.ExchangeRateService = service.NewExchangeRateService(
		container.ExternalAPIRepository,
		container.InMemoryRepository,
		cfg.Cache.MaxHistoricalDays,
	)

	// Use case layer
	container.ExchangeRateUsecase = usecase.NewExchangeRateUsecase(container.ExchangeRateService)

	// Handler layer
	container.ExchangeRateHandler = handler.NewExchangeRateHandler(container.ExchangeRateUsecase)

	return container
}
