package router

import (
	"exchange-rate-service/internal/delivery/http/middleware"
	transport "exchange-rate-service/internal/delivery/http/transport"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupRoutes(handlers transport.Handlers) *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())
	router.Use(middleware.PrometheusMetrics())

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	api := router.Group("/api/")
	{
		api.GET("/latest", gin.WrapH(handlers.Latest))
		api.GET("/convert", gin.WrapH(handlers.Convert))
		api.GET("/historical", gin.WrapH(handlers.Historical))
	}

	return router
}
