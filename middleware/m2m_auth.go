package middleware

import (
	"api-monitoring/app/m2m"
	"api-monitoring/errors"
	"api-monitoring/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func M2MAuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID := c.GetHeader("X-Client-ID")
		clientSecret := c.GetHeader("X-Client-Secret")

		if clientID == "" || clientSecret == "" {
			httpErr := errors.UnauthorizedException()
			c.AbortWithStatusJSON(httpErr.StatusCode, utils.NewHttpError("Missing client credentials", httpErr.Code, nil))
			return
		}

		var client m2m.Client
		result := db.Where("id = ? AND secret = ? AND is_active = ?", clientID, clientSecret, true).First(&client)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				httpErr := errors.UnauthorizedException()
				c.AbortWithStatusJSON(httpErr.StatusCode, utils.NewHttpError("Invalid client credentials", httpErr.Code, nil))
				return
			}
			c.AbortWithStatusJSON(500, utils.NewHttpError("Internal server error", "ERR_INTERNAL", nil))
			return
		}

		c.Set("m2m_client", &client)
		c.Next()
	}
}
