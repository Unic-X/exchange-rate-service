package di

import (
	"context"

	"exchange-rate-service/internal/delivery/http/handler"
	"exchange-rate-service/internal/domain/config"
	domain_exchange "exchange-rate-service/internal/domain/exchange"
	"exchange-rate-service/internal/infra/http_client"
	"exchange-rate-service/internal/infra/repository/api"
	"exchange-rate-service/internal/infra/repository/inmemory"
	"exchange-rate-service/internal/infra/repository/mock"
	usecase "exchange-rate-service/internal/usecase/exchange"

	"exchange-rate-service/pkg/cache"
)

type InfraContainer struct {
	HTTPClient http_client.HTTPClient
	Cache      cache.Cache
}

type RepositoryContainer struct {
	ExternalAPIRepository domain_exchange.ExchangeRateExternalRepository
	InMemoryRepository    domain_exchange.ExchangeRateCacheRepository
	MockRepository        domain_exchange.ExchangeRateExternalRepository
}

type UseCaseContainer struct {
	ExchangeRateUseCase domain_exchange.ExchangeRateUsercase
}

type HandlerContainer struct {
	ExchangeRateHandler *handler.ExchangeRateHandler
}

type AppContainer struct {
	Infra        *InfraContainer
	Repositories *RepositoryContainer
	UseCases     *UseCaseContainer
	Handlers     *HandlerContainer
	Config       *config.Config
}

func NewAppContainer(ctx context.Context, cfg *config.Config) *AppContainer {
	infra := &InfraContainer{
		HTTPClient: http_client.NewHTTPClient(cfg.FiatExternalAPI.Timeout),
		Cache:      cache.NewInMemoryCache(cfg.Cache.TTL),
	}

	fiatRepo := api.NewExternalAPIRepository(
		infra.HTTPClient,
		cfg.FiatExternalAPI.BaseURL,
		cfg.FiatExternalAPI.Secret,
	)
	cryptoRepo := api.NewCryptoAPIRepository(
		infra.HTTPClient,
		cfg.CryptoExternalAPI.BaseURL,
		cfg.CryptoExternalAPI.Secret,
	)
	mockRepository := mock.NewMockExchangeRateRepository()

	repos := &RepositoryContainer{
		ExternalAPIRepository: api.NewCompositeRepository(fiatRepo, cryptoRepo, mockRepository),
		InMemoryRepository:    inmemory.NewInMemoryRepository(infra.Cache),
		MockRepository:        mockRepository,
	}

	useCases := &UseCaseContainer{
		ExchangeRateUseCase: usecase.NewExchangeRateUseCase(
			repos.ExternalAPIRepository,
			repos.InMemoryRepository,
			cfg.Cache.MaxHistoricalDays,
		),
	}

	handlers := &HandlerContainer{
		ExchangeRateHandler: handler.NewExchangeRateHandler(useCases.ExchangeRateUseCase),
	}

	app := &AppContainer{
		Infra:        infra,
		Repositories: repos,
		UseCases:     useCases,
		Handlers:     handlers,
		Config:       cfg,
	}

	go usecase.StartRateRefreshTicker(useCases.ExchangeRateUseCase, ctx, cfg.Cache.RefreshInterval)

	return app
}
