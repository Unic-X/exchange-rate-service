package handler

import (
	"net/http"
	"strconv"
	"time"

	domain_exchange "exchange-rate-service/internal/domain/exchange"

	"github.com/gin-gonic/gin"
)

type ExchangeRateHandler struct {
	usecase domain_exchange.ExchangeRateUsercase
}

func NewExchangeRateHandler(u domain_exchange.ExchangeRateUsercase) *ExchangeRateHandler {
	return &ExchangeRateHandler{usecase: u}
}

func ParseDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, nil
	}
	targetDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return targetDate, nil
}

func (h *ExchangeRateHandler) ConvertAmount(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	amountStr := c.Query("amount")

	fromDate := c.Query("fromDate")
	toDate := c.Query("toDate")

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
		return
	}

	fromTargetDate, err := ParseDate(fromDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid from date format. Use YYYY-MM-DD"})
		return
	}

	toTargetDate, err := ParseDate(toDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid to date format. Use YYYY-MM-DD"})
		return
	}

	fromAmount, toAmount, fromRate, toRate, err := h.usecase.ConvertAmount(c, from, to, amount, fromTargetDate, toTargetDate)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"from":              from,
		"to":                to,
		"original_amount":   amount,
		"converted_at_from": fromAmount,
		"converted_at_to":   toAmount,
		"from_rate":         fromRate,
		"to_rate":           toRate,
		"from_date":         fromTargetDate,
		"to_date":           toTargetDate,
	})
}
