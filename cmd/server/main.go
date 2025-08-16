package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"exchange-rate-service/internal/config"
	"exchange-rate-service/internal/delivery/http/router"
	"exchange-rate-service/internal/di"
	"exchange-rate-service/pkg/logger"
)

func main() {
	cfg := config.Load()

	logger.Info("Starting Exchange Rate Service...")
	logger.Infof("Server will start on %s:%s", cfg.Server.Host, cfg.Server.Port)

	container := di.NewContainer(cfg)

	r := router.SetupRoutes(container.HTTPHandlers)

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go startRateRefreshTicker(ctx, container, cfg.Cache.RefreshInterval)

	go func() {
		logger.Infof("Server starting on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("Server failed to start: %v", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Cancel rate refresh ticker
	cancel()

	// Shutdown server gracefully
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Errorf("Server forced to shutdown: %v", err)
	} else {
		logger.Info("Server exited gracefully")
	}

	logger.Sync()
}

func startRateRefreshTicker(ctx context.Context, container *di.Container, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	logger.Infof("Starting rate refresh ticker with interval: %v", interval)

	// Initial refresh on startup
	if err := container.ExchangeRateService.RefreshRates(ctx); err != nil {
		logger.Errorf("Initial rate refresh failed: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			logger.Info("Rate refresh ticker stopped")
			return
		case <-ticker.C:
			logger.Info("Refreshing exchange rates...")
			if err := container.ExchangeRateService.RefreshRates(ctx); err != nil {
				logger.Errorf("Rate refresh failed: %v", err)
			} else {
				logger.Info("Exchange rates refreshed successfully")
			}
		}
	}
}
