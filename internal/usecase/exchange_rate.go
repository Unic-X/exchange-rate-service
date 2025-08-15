package usecase

import (
	"context"
	"time"

	"exchange-rate-service/internal/domain/service"
)

type ExchangeRateUsecase interface {
	GetLatestRate(from, to string) (float64, error)
	ConvertAmount(from, to string, amount float64, date time.Time) (float64, float64, error)
	GetHistoricalRate(from, to string, date time.Time) (float64, error)
}

type exchangeRateUsecase struct {
	service service.ExchangeRateService
}

func NewExchangeRateUsecase(service service.ExchangeRateService) ExchangeRateUsecase {
	return &exchangeRateUsecase{
		service: service,
	}
}

func (u *exchangeRateUsecase) GetLatestRate(from, to string) (float64, error) {
	ctx := context.Background()
	return u.service.GetLatestRate(ctx, from, to)
}

func (u *exchangeRateUsecase) ConvertAmount(from, to string, amount float64, date time.Time) (float64, float64, error) {
	ctx := context.Background()
	return u.service.ConvertAmount(ctx, from, to, amount, date)
}

func (u *exchangeRateUsecase) GetHistoricalRate(from, to string, date time.Time) (float64, error) {
	ctx := context.Background()
	return u.service.GetHistoricalRate(ctx, from, to, date)
}
