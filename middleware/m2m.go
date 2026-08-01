package middleware

import (
	"api-monitoring/errors"
	"api-monitoring/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func M2MAPIKeyMiddleware(validAPIKeys []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			httpErr := errors.UnauthorizedException()
			c.AbortWithStatusJSON(httpErr.StatusCode, utils.NewHttpError("Missing API key", httpErr.Code, nil))
			return
		}

		valid := false
		for _, key := range validAPIKeys {
			if apiKey == key {
				valid = true
				break
			}
		}

		if !valid {
			httpErr := errors.UnauthorizedException()
			c.AbortWithStatusJSON(httpErr.StatusCode, utils.NewHttpError("Invalid API key", httpErr.Code, nil))
			return
		}

		c.Next()
	}
}

func M2MIPWhitelistMiddleware(allowedIPs []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(allowedIPs) == 0 {
			c.Next()
			return
		}

		clientIP := c.ClientIP()

		allowed := false
		for _, ip := range allowedIPs {
			if clientIP == ip {
				allowed = true
				break
			}
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, utils.NewHttpError("IP address not whitelisted", "ERR_FORBIDDEN", nil))
			return
		}

		c.Next()
	}
}
