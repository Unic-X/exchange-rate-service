package exchange

import (
	"context"
	"time"

	domain_exchange "exchange-rate-service/internal/domain/exchange"
	"exchange-rate-service/pkg/logger"
)

func StartRateRefreshTicker(useCase domain_exchange.ExchangeRateUsercase, ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	logger.Infof("Starting rate refresh ticker with interval: %v", interval)

	if err := useCase.RefreshRates(ctx); err != nil {
		logger.Errorf("Initial rate refresh failed: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			logger.Info("Rate refresh ticker stopped")
			return
		case <-ticker.C:
			logger.Info("Refreshing exchange rates...")
			if err := useCase.RefreshRates(ctx); err != nil {
				logger.Errorf("Rate refresh failed: %v", err)
			} else {
				logger.Info("Exchange rates refreshed successfully")
			}
		}
	}
}
