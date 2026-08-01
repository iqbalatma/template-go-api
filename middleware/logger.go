package middleware

import (
	"api-monitoring/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const loggerKey = "logger"

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()
		c.Set("RequestID", requestID)

		entry := config.AppLogger.WithFields(logrus.Fields{
			"request_id": requestID,
			"ip":         c.ClientIP(),
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
		})
		c.Set(loggerKey, entry)

		start := time.Now()
		c.Next()

		if c.Request.Method == http.MethodOptions {
			return
		}

		entry.WithFields(logrus.Fields{
			"status":  c.Writer.Status(),
			"latency": time.Since(start).String(),
		}).Info("request completed")
	}
}

func Logger(c *gin.Context) *logrus.Entry {
	if entry, ok := c.Get(loggerKey); ok {
		return entry.(*logrus.Entry)
	}
	return config.AppLogger.WithField("ip", c.ClientIP())
}
