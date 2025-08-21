package handler

import (
	"net/http"
	"strconv"
	"time"

	"exchange-rate-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

type ExchangeRateHandler struct {
	usecase usecase.ExchangeRateUsecase
}

func NewExchangeRateHandler(u usecase.ExchangeRateUsecase) *ExchangeRateHandler {
	return &ExchangeRateHandler{usecase: u}
}

func ParseDate(dateStr string) (time.Time, error) {
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

	fromRate, toRate, converted, err := h.usecase.ConvertAmount(c, from, to, amount, fromTargetDate, toTargetDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"from":             from,
		"to":               to,
		"original_amount":  amount,
		"converted_amount": converted,
		"from_rate":        fromRate,
		"to_rate":          toRate,
		"difference":       toRate - fromRate,
	})
}
