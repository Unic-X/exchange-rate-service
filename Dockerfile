FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
    
COPY . .

ARG BUILD_MODE=debug
ARG VERSION=dev

RUN if [ "$BUILD_MODE" = "release" ]; then \
        CGO_ENABLED=0 GOOS=linux go build \
        -a -installsuffix cgo \
        -ldflags="-w -s -X main.version=$VERSION" \
        -o exchange-rate-service ./cmd/server; \
    else \
        go build -o exchange-rate-service ./cmd/server; \
    fi
    

FROM alpine:latest
    
WORKDIR /app
RUN apk --no-cache add ca-certificates
    
COPY --from=builder /app/exchange-rate-service .
COPY --from=builder /app/.env* ./
    
EXPOSE 8080
CMD ["./exchange-rate-service"]    