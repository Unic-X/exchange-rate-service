# Exchange Rate Service

An exchange rate service built with Go, implementing Clean Architecture principles with Dependency Injection. This service provides real-time and historical exchange rate data with caching capabilities.


### Postman Collection

Import the Postman collection to test the API endpoints:

**[📋 Postman Collection](https://arman22102-3102413.postman.co/workspace/Arman-Singh-Kshatri's-Workspace~ed568349-d9d4-459b-89de-aa9ad2b35f81/collection/47632237-1c467d05-c96c-4a6d-86c0-98dc8536e356?action=share&creator=47632237)**


## 🏗️ Architecture Overview

This project follows **Clean Architecture** principles with clear separation of concerns:

```
├── cmd/server/           # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── delivery/http/   # HTTP handlers, middleware, routes
│   ├── di/              # Dependency Injection container
│   ├── domain/          # Business entities, interfaces
│   ├── infra/           # External dependencies (HTTP client, cache)
│   └── usecase/         # Business logic layer
└── pkg/                 # Shared utilities (logger)
```

### Clean Architecture Layers

1. **Domain Layer** (`internal/domain/`): Core business entities and interfaces
2. **Use Case Layer** (`internal/usecase/`): Business logic orchestration
3. **Infrastructure Layer** (`internal/infra/`): External dependencies and data sources
4. **Delivery Layer** (`internal/delivery/`): HTTP handlers and middleware
5. **Dependency Injection** (`internal/di/`): Wires all dependencies together

### Dependency Injection Implementation

The DI container (`internal/di/dependency_injection.go`) manages all dependencies:

```go
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
```

Dependencies are wired from the bottom up:
- **Infrastructure** → **Repository** → **Service** → **Use Case** → **Handler**

## 🚀 Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.24+ (for local development)

### Using Docker Compose (Recommended)

1. **Clone the repository**
```bash
git clone <repository-url>
cd ExchangeRateService
```

2. **Create environment file**
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. **Build and run with Docker Compose**
```bash
docker-compose up --build
```

The service will be available at `http://localhost:8080`

### Using Docker Build

1. **Build the Docker image**
```bash
docker build -t exchange-rate-service .
```

2. **Run the container**
```bash
docker run -p 8080:8080 --env-file .env exchange-rate-service
```

### Local Development

1. **Install dependencies**
```bash
go mod download
```

2. **Run the service**
```bash
go run cmd/server/main.go
```

## 🔧 Environment Configuration

Create a `.env` file in the root directory with the following structure:

### Server Configuration
```env
# Server settings
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_READ_TIMEOUT=30s
SERVER_WRITE_TIMEOUT=30s
```

### External API Configuration
```env
# External exchange rate API settings
EXTERNAL_API_BASE_URL=https://v6.exchangerate-api.com/v6
EXTERNAL_API_SECRET=your_api_key_here
EXTERNAL_API_TIMEOUT=10s
EXTERNAL_API_RETRY_ATTEMPTS=3
EXTERNAL_API_RETRY_DELAY=1s
```

### Cache Configuration
```env
# Caching settings
CACHE_TTL=1h
CACHE_REFRESH_INTERVAL=1h
MAX_HISTORICAL_DAYS=90
```

### Environment Variables Reference

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `SERVER_HOST` | Server bind address | `0.0.0.0` | No |
| `SERVER_PORT` | Server port | `8080` | No |
| `SERVER_READ_TIMEOUT` | HTTP read timeout | `30s` | No |
| `SERVER_WRITE_TIMEOUT` | HTTP write timeout | `30s` | No |
| `EXTERNAL_API_BASE_URL` | External API base URL | `https://v6.exchangerate-api.com/v6` | No |
| `EXTERNAL_API_SECRET` | External API key | `secret` | **Yes** |
| `EXTERNAL_API_TIMEOUT` | API request timeout | `10s` | No |
| `EXTERNAL_API_RETRY_ATTEMPTS` | Number of retry attempts | `3` | No |
| `EXTERNAL_API_RETRY_DELAY` | Delay between retries | `1s` | No |
| `CACHE_TTL` | Cache time-to-live | `1h` | No |
| `CACHE_REFRESH_INTERVAL` | Cache refresh interval | `1h` | No |
| `MAX_HISTORICAL_DAYS` | Maximum historical data days | `90` | No |

## 📡 API Documentation

### Available Endpoints

- `GET /health` - Health check endpoint
- `GET /api/v1/exchange-rate/{from}/{to}` - Get current exchange rate
- `GET /api/v1/exchange-rate/{from}/{to}/historical` - Get historical exchange rates

## 🛠️ Development

### Project Structure

```
exchange-rate-service/
├── cmd/server/main.go              # Application entry point
├── internal/
│   ├── config/config.go            # Configuration management
│   ├── delivery/http/
│   │   ├── handler/                # HTTP request handlers
│   │   ├── middleware/             # HTTP middleware (logging, etc.)
│   │   └── routes/                 # Route definitions
│   ├── di/dependency_injection.go  # Dependency injection container
│   ├── domain/
│   │   ├── entity/                 # Business entities
│   │   ├── repository/             # Repository interfaces
│   │   └── service/                # Service interfaces
│   ├── infra/
│   │   ├── cache/                  # Cache implementations
│   │   ├── http_client/            # HTTP client wrapper
│   │   └── repository/             # Repository implementations
│   └── usecase/                    # Business logic use cases
├── pkg/logger/                     # Shared logging utility
├── docker-compose.yml              # Docker Compose configuration
├── Dockerfile                      # Docker build configuration
├── go.mod                          # Go module dependencies
└── .env                           # Environment variables
```

### Key Features

- **Clean Architecture**: Separation of concerns with clear layer boundaries
- **Dependency Injection**: Centralized dependency management
- **Caching**: In-memory caching with configurable TTL
- **Retry Logic**: Configurable retry mechanism for external API calls
- **Structured Logging**: JSON-structured logging with different levels
- **Health Checks**: Built-in health check endpoint
- **Docker Support**: Multi-stage Docker build for optimized images
- **Configuration Management**: Environment-based configuration

### Building for Production

The Dockerfile uses multi-stage builds for optimized production images:

1. **Builder stage**: Compiles the Go application
2. **Runtime stage**: Minimal Alpine Linux image with the compiled binary

This results in a lightweight production image (~15MB) with only the necessary runtime dependencies.

## 📝 Logging

The service uses structured JSON logging with different levels:
- `INFO`: General information and successful operations
- `ERROR`: Error conditions and failures
- `DEBUG`: Detailed debugging information

Logs include request details, response times, and error messages for comprehensive monitoring.

## 🔒 Security

- API keys are managed through environment variables
- HTTP timeouts prevent hanging requests
- Input validation on all endpoints
- Structured error responses without sensitive information exposure

## 🚀 Deployment

The service is containerized and ready for deployment on any container orchestration platform:

- **Docker**: Single container deployment
- **Docker Compose**: Local development and testing
- **Kubernetes**: Production-ready with proper resource limits

## 📊 Monitoring

The service includes:
- Structured logging(currently only STDOUT) for log aggregation systems
- Request/response timing metrics
- Error tracking and reporting

---

**Built with ❤️ using Go and Clean Architecture principles**
