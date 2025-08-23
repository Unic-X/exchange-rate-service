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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	container := di.NewAppContainer(ctx, cfg)

	r := router.SetupRoutes(container.Handlers.ExchangeRateHandler)

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

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

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Errorf("Server forced to shutdown: %v", err)
	} else {
		logger.Info("Server exited gracefully")
	}

	logger.Sync()
}
