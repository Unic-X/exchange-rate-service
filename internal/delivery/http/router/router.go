package router

import (
	"exchange-rate-service/internal/delivery/http/handler"
	"exchange-rate-service/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupRoutes(exchangeRateHandler *handler.ExchangeRateHandler) *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())
	router.Use(middleware.PrometheusMetrics())

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	api := router.Group("/api/")
	{
		api.GET("/convert", exchangeRateHandler.ConvertAmount)

	}

	return router
}
