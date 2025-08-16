FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
    
COPY . .
RUN go build -o exchange-rate-service ./cmd/server
    
# Stage 2 

FROM alpine:latest
    
WORKDIR /app
RUN apk --no-cache add ca-certificates
    
COPY --from=builder /app/exchange-rate-service .
    
EXPOSE 8080
CMD ["./exchange-rate-service"]    