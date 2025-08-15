package handler

import (
	"net/http"
	"strconv"
	"time"

	"exchange-rate-service/pkg/logger"

	"github.com/gin-gonic/gin"
)

type ExchangeRateHandler struct {
	// TODO : add the use case dependency later
}

func NewExchangeRateHandler() *ExchangeRateHandler {
	return &ExchangeRateHandler{}
}

func (h *ExchangeRateHandler) GetLatestRate(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")

	if from == "" || to == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Both 'from' and 'to' currency parameters are required",
		})
		return
	}

	// TODO: Implement actual logic
	logger.Infof("Getting latest rate from %s to %s", from, to)

	c.JSON(http.StatusOK, gin.H{
		"from": from,
		"to":   to,
		"rate": 83.25,
		"date": time.Now().Format("2006-01-02"),
	})
}

func (h *ExchangeRateHandler) ConvertAmount(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	amountStr := c.Query("amount")
	dateStr := c.Query("date")

	if from == "" || to == "" || amountStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Parameters 'from', 'to', and 'amount' are required",
		})
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid amount format",
		})
		return
	}

	var targetDate time.Time
	if dateStr != "" {
		targetDate, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid date format. Use YYYY-MM-DD",
			})
			return
		}
	} else {
		targetDate = time.Now()
	}

	// TODO: Implement actual conversion logic
	logger.Infof("Converting %.2f %s to %s for date %s", amount, from, to, targetDate.Format("2006-01-02"))

	// temporary (1 USD = 83.25 INR)
	convertedAmount := amount * 83.25

	c.JSON(http.StatusOK, gin.H{
		"from":             from,
		"to":               to,
		"original_amount":  amount,
		"converted_amount": convertedAmount,
		"rate":             83.25,
		"date":             targetDate.Format("2006-01-02"),
	})
}

func (h *ExchangeRateHandler) GetHistoricalRate(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	dateStr := c.Query("date")

	if from == "" || to == "" || dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Parameters 'from', 'to', and 'date' are required",
		})
		return
	}

	targetDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid date format. Use YYYY-MM-DD",
		})
		return
	}

	maxPastDate := time.Now().AddDate(0, 0, -90)
	if targetDate.Before(maxPastDate) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Date cannot be older than 90 days",
		})
		return
	}

	if targetDate.After(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Date cannot be in the future",
		})
		return
	}

	// TODO: Implement actual historical rate logic
	logger.Infof("Getting historical rate from %s to %s for date %s", from, to, targetDate.Format("2006-01-02"))

	c.JSON(http.StatusOK, gin.H{
		"from": from,
		"to":   to,
		"rate": 82.85,
		"date": targetDate.Format("2006-01-02"),
	})
}
