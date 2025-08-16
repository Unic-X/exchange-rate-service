package transport

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"exchange-rate-service/internal/delivery/endpoint"
	"exchange-rate-service/pkg/logger"

	kittransport "github.com/go-kit/kit/transport/http"
)

type Handlers struct {
	Latest     http.Handler
	Historical http.Handler
	Convert    http.Handler
}

func NewHandlers(eps endpoint.Endpoints) Handlers {
	return Handlers{
		Latest:     kittransport.NewServer(eps.GetLatestRate, decodeGetLatestRateRequest, encodeResponse, kittransport.ServerErrorEncoder(errorEncoder)),
		Historical: kittransport.NewServer(eps.GetHistoricalRate, decodeGetHistoricalRateRequest, encodeResponse, kittransport.ServerErrorEncoder(errorEncoder)),
		Convert:    kittransport.NewServer(eps.ConvertAmount, decodeConvertAmountRequest, encodeResponse, kittransport.ServerErrorEncoder(errorEncoder)),
	}
}

func decodeGetLatestRateRequest(_ context.Context, r *http.Request) (any, error) {
	q := r.URL.Query()
	from := q.Get("from")
	to := q.Get("to")
	return endpoint.GetLatestRateRequest{From: from, To: to}, nil
}

func decodeGetHistoricalRateRequest(_ context.Context, r *http.Request) (any, error) {
	q := r.URL.Query()
	from := q.Get("from")
	to := q.Get("to")
	dateStr := q.Get("date")
	var date time.Time
	var err error
	if dateStr == "" {
		return nil, httpError{Status: http.StatusBadRequest, Msg: "Invalid date format. Use YYYY-MM-DD"}
	}
	date, err = time.Parse("2006-01-02", dateStr)
	if err != nil {
		logger.Warnf("Invalid date format in request: %s", dateStr)
		return nil, httpError{Status: http.StatusBadRequest, Msg: "Invalid date format. Use YYYY-MM-DD"}
	}
	return endpoint.GetHistoricalRateRequest{From: from, To: to, Date: date}, nil
}

func decodeConvertAmountRequest(_ context.Context, r *http.Request) (any, error) {
	q := r.URL.Query()
	from := q.Get("from")
	to := q.Get("to")
	amountStr := q.Get("amount")
	dateStr := q.Get("date")

	var amount float64
	if amountStr != "" {
		if v, err := strconv.ParseFloat(amountStr, 64); err == nil {
			amount = v
		} else {
			return nil, httpError{Status: http.StatusBadRequest, Msg: "Invalid amount"}
		}
	} else {
		return nil, httpError{Status: http.StatusBadRequest, Msg: "Invalid amount"}
	}

	var date time.Time
	if dateStr != "" {
		if v, err := time.Parse("2006-01-02", dateStr); err == nil {
			date = v
		} else {
			logger.Warnf("Invalid date format in request: %s", dateStr)
			return nil, httpError{Status: http.StatusBadRequest, Msg: "Invalid date format. Use YYYY-MM-DD"}
		}
	} else {
		date = time.Now()
	}

	return endpoint.ConvertAmountRequest{From: from, To: to, Amount: amount, Date: date}, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response any) error {
	_ = _contextWithContentType(w)
	return kittransport.EncodeJSONResponse(context.Background(), w, response)
}

func _contextWithContentType(w http.ResponseWriter) context.Context {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return context.Background()
}

type httpError struct {
	Status int
	Msg    string
}

func (e httpError) Error() string { return e.Msg }

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	code := http.StatusBadRequest
	if he, ok := err.(httpError); ok {
		code = he.Status
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = kittransport.EncodeJSONResponse(context.Background(), w, map[string]string{"error": err.Error()})
}
