package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server      ServerConfig
	ExternalAPI ExternalAPIConfig
	Cache       CacheConfig
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

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 30*time.Second),
		},
		ExternalAPI: ExternalAPIConfig{
			BaseURL:       getEnv("EXTERNAL_API_BASE_URL", "https://v6.exchangerate-api.com/v6"),
			Secret:        getEnv("EXTERNAL_API_SECRET", "secret"),
			Timeout:       getDurationEnv("EXTERNAL_API_TIMEOUT", 10*time.Second),
			RetryAttempts: getIntEnv("EXTERNAL_API_RETRY_ATTEMPTS", 3),
			RetryDelay:    getDurationEnv("EXTERNAL_API_RETRY_DELAY", 1*time.Second),
		},
		Cache: CacheConfig{
			TTL:               getDurationEnv("CACHE_TTL", 1*time.Hour),
			RefreshInterval:   getDurationEnv("CACHE_REFRESH_INTERVAL", 1*time.Hour),
			MaxHistoricalDays: getIntEnv("MAX_HISTORICAL_DAYS", 90),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
