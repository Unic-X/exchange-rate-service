package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type mockUsecase struct {
	latest float64
	hist   float64
	rate   float64
	amt    float64
	err    error
}

func (m *mockUsecase) GetLatestRate(from, to string) (float64, error) { return m.latest, m.err }
func (m *mockUsecase) ConvertAmount(from, to string, amount float64, date time.Time) (float64, float64, error) {
	return m.rate, m.amt, m.err
}
func (m *mockUsecase) GetHistoricalRate(from, to string, date time.Time) (float64, error) {
	return m.hist, m.err
}

func TestGetLatestRate_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mu := &mockUsecase{latest: 1.5}
	h := NewExchangeRateHandler(mu)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	q := url.Values{"from": {"EUR"}, "to": {"USD"}}
	req := httptest.NewRequest(http.MethodGet, "/api/latest?"+q.Encode(), nil)
	c.Request = req

	h.ConvertAmount(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestGetLatestRate_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mu := &mockUsecase{err: errors.New("this is the error that would propagate")}
	h := NewExchangeRateHandler(mu)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	q := url.Values{"from": {"EUR"}, "to": {"USD"}}
	c.Request = httptest.NewRequest(http.MethodGet, "/api/latest?"+q.Encode(), nil)

	h.ConvertAmount(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestConvertAmount_InvalidAmount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewExchangeRateHandler(&mockUsecase{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	q := url.Values{"from": {"EUR"}, "to": {"USD"}, "amount": {"-1"}}
	c.Request = httptest.NewRequest(http.MethodGet, "/api/convert?"+q.Encode(), nil)

	h.ConvertAmount(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestConvertAmount_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mu := &mockUsecase{rate: 2, amt: 20}
	h := NewExchangeRateHandler(mu)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	q := url.Values{"from": {"EUR"}, "to": {"USD"}, "amount": {"10"}, "date": {"2024-01-02"}}
	c.Request = httptest.NewRequest(http.MethodGet, "/api/convert?"+q.Encode(), nil)

	h.ConvertAmount(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestGetHistoricalRate_InvalidDate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewExchangeRateHandler(&mockUsecase{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	q := url.Values{"from": {"EUR"}, "to": {"USD"}, "date": {"bad"}}
	c.Request = httptest.NewRequest(http.MethodGet, "/api/historical?"+q.Encode(), nil)

	h.ConvertAmount(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestGetHistoricalRate_TooOld(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewExchangeRateHandler(&mockUsecase{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	old := time.Now().AddDate(0, 0, -91).Format("2006-01-02")
	q := url.Values{"from": {"EUR"}, "to": {"USD"}, "date": {old}}
	c.Request = httptest.NewRequest(http.MethodGet, "/api/historical?"+q.Encode(), nil)

	h.ConvertAmount(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestGetHistoricalRate_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mu := &mockUsecase{hist: 1.1}
	h := NewExchangeRateHandler(mu)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	date := time.Now().AddDate(0, 0, -5).Format("2006-01-02")
	q := url.Values{"from": {"EUR"}, "to": {"USD"}, "date": {date}}
	c.Request = httptest.NewRequest(http.MethodGet, "/api/historical?"+q.Encode(), nil)

	h.ConvertAmount(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
}
