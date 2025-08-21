package usecase

import (
	"context"
	"time"

	"exchange-rate-service/internal/domain/service"
)

type ExchangeRateUsecase interface {
	ConvertAmount(ctx context.Context, from, to string, amount float64, fromDate, toDate time.Time) (float64, float64, float64, error)
}

type exchangeRateUsecase struct {
	service service.ExchangeRateService
}

func NewExchangeRateUsecase(service service.ExchangeRateService) ExchangeRateUsecase {
	return &exchangeRateUsecase{
		service: service,
	}
}

func (u *exchangeRateUsecase) ConvertAmount(ctx context.Context, from, to string, amount float64, fromDate, toDate time.Time) (float64, float64, float64, error) {
	return u.service.ConvertAmount(ctx, from, to, amount, fromDate, toDate)
}
