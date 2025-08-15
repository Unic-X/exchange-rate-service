package router

import (
	"exchange-rate-service/internal/delivery/http/handler"
	"exchange-rate-service/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(exchangeRateHandler *handler.ExchangeRateHandler) *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "exchange-rate-service",
		})
	})

	api := router.Group("/api/")
	{
		api.GET("/latest", exchangeRateHandler.GetLatestRate)
		api.GET("/convert", exchangeRateHandler.ConvertAmount)
		api.GET("/historical", exchangeRateHandler.GetHistoricalRate)
	}

	return router
}
