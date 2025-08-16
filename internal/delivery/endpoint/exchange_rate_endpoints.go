package endpoint

import (
	"context"
	"time"

	"exchange-rate-service/internal/usecase"

	kitendpoint "github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GetLatestRate     kitendpoint.Endpoint
	GetHistoricalRate kitendpoint.Endpoint
	ConvertAmount     kitendpoint.Endpoint
}

type GetLatestRateRequest struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type GetLatestRateResponse struct {
	From  string  `json:"from"`
	To    string  `json:"to"`
	Rate  float64 `json:"rate"`
	Date  string  `json:"date"`
	Error string  `json:"error,omitempty"`
}

type GetHistoricalRateRequest struct {
	From string    `json:"from"`
	To   string    `json:"to"`
	Date time.Time `json:"date"`
}

type GetHistoricalRateResponse struct {
	From  string  `json:"from"`
	To    string  `json:"to"`
	Rate  float64 `json:"rate"`
	Date  string  `json:"date"`
	Error string  `json:"error,omitempty"`
}

type ConvertAmountRequest struct {
	From   string    `json:"from"`
	To     string    `json:"to"`
	Amount float64   `json:"amount"`
	Date   time.Time `json:"date"`
}

type ConvertAmountResponse struct {
	From            string  `json:"from"`
	To              string  `json:"to"`
	OriginalAmount  float64 `json:"original_amount"`
	ConvertedAmount float64 `json:"converted_amount"`
	Rate            float64 `json:"rate"`
	Date            string  `json:"date"`
	Error           string  `json:"error,omitempty"`
}

func MakeEndpoints(uc usecase.ExchangeRateUsecase) Endpoints {
	return Endpoints{
		GetLatestRate:     makeGetLatestRateEndpoint(uc),
		GetHistoricalRate: makeGetHistoricalRateEndpoint(uc),
		ConvertAmount:     makeConvertAmountEndpoint(uc),
	}
}

func makeGetLatestRateEndpoint(uc usecase.ExchangeRateUsecase) kitendpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(GetLatestRateRequest)
		rate, err := uc.GetLatestRate(req.From, req.To)
		if err != nil {
			return GetLatestRateResponse{From: req.From, To: req.To, Error: err.Error()}, nil
		}
		return GetLatestRateResponse{From: req.From, To: req.To, Rate: rate, Date: time.Now().Format("2006-01-02")}, nil
	}
}

func makeGetHistoricalRateEndpoint(uc usecase.ExchangeRateUsecase) kitendpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(GetHistoricalRateRequest)
		rate, err := uc.GetHistoricalRate(req.From, req.To, req.Date)
		if err != nil {
			return GetHistoricalRateResponse{From: req.From, To: req.To, Date: req.Date.Format("2006-01-02"), Error: err.Error()}, nil
		}
		return GetHistoricalRateResponse{From: req.From, To: req.To, Rate: rate, Date: req.Date.Format("2006-01-02")}, nil
	}
}

func makeConvertAmountEndpoint(uc usecase.ExchangeRateUsecase) kitendpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(ConvertAmountRequest)
		rate, converted, err := uc.ConvertAmount(req.From, req.To, req.Amount, req.Date)
		if err != nil {
			return ConvertAmountResponse{From: req.From, To: req.To, OriginalAmount: req.Amount, Date: req.Date.Format("2006-01-02"), Error: err.Error()}, nil
		}
		return ConvertAmountResponse{From: req.From, To: req.To, OriginalAmount: req.Amount, ConvertedAmount: converted, Rate: rate, Date: req.Date.Format("2006-01-02")}, nil
	}
}
