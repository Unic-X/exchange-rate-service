package middleware

import (
	"exchange-rate-service/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logMessage := "[%s] %s %s %d %s %s"
		logArgs := []any{
			param.TimeStamp.Format(time.RFC3339),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.ClientIP,
		}

		if param.StatusCode >= 400 {
			logger.Errorf(logMessage, logArgs...)
		} else {
			logger.Infof(logMessage, logArgs...)
		}
		return ""
	})
}
