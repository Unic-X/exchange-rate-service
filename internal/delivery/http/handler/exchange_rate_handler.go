package handler

import (
	"net/http"
	"strconv"
	"time"

	"exchange-rate-service/internal/usecase"
	"exchange-rate-service/pkg/logger"

	"github.com/gin-gonic/gin"
)

type ExchangeRateHandler struct {
	usecase usecase.ExchangeRateUsecase
}

func NewExchangeRateHandler(u usecase.ExchangeRateUsecase) *ExchangeRateHandler {
	return &ExchangeRateHandler{usecase: u}
}

func (h *ExchangeRateHandler) GetLatestRate(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")

	rate, err := h.usecase.GetLatestRate(from, to)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Infof("Getting latest rate from %s to %s", from, to)

	c.JSON(http.StatusOK, gin.H{
		"from": from,
		"to":   to,
		"rate": rate,
		"date": time.Now().Format("2006-01-02"),
	})
}

func (h *ExchangeRateHandler) ConvertAmount(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	amountStr := c.Query("amount")
	dateStr := c.Query("date")

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
		return
	}

	var targetDate time.Time
	if dateStr != "" {
		targetDate, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
			return
		}
	} else {
		targetDate = time.Now()
	}

	rate, converted, err := h.usecase.ConvertAmount(from, to, amount, targetDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Infof("Converting %.2f %s to %s for date %s", amount, from, to, targetDate.Format("2006-01-02"))

	c.JSON(http.StatusOK, gin.H{
		"from":             from,
		"to":               to,
		"original_amount":  amount,
		"converted_amount": converted,
		"rate":             rate,
		"date":             targetDate.Format("2006-01-02"),
	})
}

func (h *ExchangeRateHandler) GetHistoricalRate(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	dateStr := c.Query("date")

	targetDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	maxPastDate := time.Now().AddDate(0, 0, -90)
	if targetDate.Before(maxPastDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date cannot be older than 90 days"})
		return
	}

	if targetDate.After(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date cannot be in the future"})
		return
	}

	rate, err := h.usecase.GetHistoricalRate(from, to, targetDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Infof("Getting historical rate from %s to %s for date %s", from, to, targetDate.Format("2006-01-02"))

	c.JSON(http.StatusOK, gin.H{
		"from": from,
		"to":   to,
		"rate": rate,
		"date": targetDate.Format("2006-01-02"),
	})
}
