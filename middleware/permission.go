package middleware

import (
	user2 "api-monitoring/app/user"
	"api-monitoring/enums"
	errors2 "api-monitoring/errors"
	"api-monitoring/utils"

	"github.com/gin-gonic/gin"
)

func RequirePermission(permission enums.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, exists := c.Get("user")
		if !exists {
			httpErr := errors2.UnauthorizedException()
			c.AbortWithStatusJSON(httpErr.StatusCode, utils.NewHttpError(httpErr.Message, httpErr.Code, nil))
			return
		}

		user := raw.(*user2.User)
		for _, role := range user.Roles {
			if role.Name == string(enums.SUPERADMIN) {
				c.Next()
				return
			}
		}

		for _, role := range user.Roles {
			for _, perm := range role.Permissions {
				if perm.Name == string(permission) {
					c.Next()
					return
				}
			}
		}

		httpErr := errors2.InvalidAction("You do not have permission to perform this action")
		c.AbortWithStatusJSON(httpErr.StatusCode, utils.NewHttpError(httpErr.Message, httpErr.Code, nil))
	}
}
