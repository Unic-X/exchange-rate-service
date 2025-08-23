package config

import "time"

type Config struct {
	Server            ServerConfig
	FiatExternalAPI   ExternalAPIConfig
	CryptoExternalAPI ExternalAPIConfig
	Cache             CacheConfig
}

type ServerConfig struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type ExternalAPIConfig struct {
	BaseURL       string
	Timeout       time.Duration
	Secret        string
	RetryAttempts int
	RetryDelay    time.Duration
}

type CacheConfig struct {
	TTL               time.Duration
	RefreshInterval   time.Duration
	MaxHistoricalDays int
}
