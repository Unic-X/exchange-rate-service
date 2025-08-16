package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ExchangeRateRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "exchange_rate_requests_total",
			Help: "Total number of exchange rate requests",
		},
		[]string{"from_currency", "to_currency", "request_type"},
	)

	ExchangeRateRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "exchange_rate_request_duration_seconds",
			Help:    "Duration of exchange rate requests in seconds",
			Buckets: []float64{0.001, 0.01, 0.1, 0.5, 1.0, 2.5, 5.0, 10.0},
		},
		[]string{"from_currency", "to_currency", "request_type"},
	)

	CacheHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
		[]string{"cache_type"},
	)

	CacheMisses = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
		[]string{"cache_type"},
	)

	ExternalAPIRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "external_api_requests_total",
			Help: "Total number of external API requests",
		},
		[]string{"api_endpoint", "status"},
	)

	ExternalAPIRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "external_api_request_duration_seconds",
			Help:    "Duration of external API requests in seconds",
			Buckets: []float64{0.1, 0.5, 1.0, 2.0, 5.0, 10.0, 30.0},
		},
		[]string{"api_endpoint", "status"},
	)

	ActiveConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_connections",
			Help: "Number of active connections",
		},
	)

	CacheSize = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "cache_size",
			Help: "Current size of the cache",
		},
	)
)
