package middleware

import (
	user2 "api-monitoring/app/user"
	"api-monitoring/config"
	errors2 "api-monitoring/errors"
	"api-monitoring/utils"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/iqbalatma/gofortify"
	"gorm.io/gorm"
)

func RefreshMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token, err = c.Cookie("refresh_token")
		if err != nil {
			httpErr := errors2.InvalidTokenTypeException()
			c.AbortWithStatusJSON(httpErr.StatusCode, utils.NewHttpError(httpErr.Message, httpErr.Code, nil))
			return
		}
		payload, err := gofortify.ValidateRefreshToken(
			&token,
		)

		if err != nil {
			var httpErr *utils.HTTPError

			switch {
			case errors.Is(err, gofortify.ErrInvalidTokenType):
				httpErr = errors2.InvalidTokenTypeException()
			}

			if httpErr == nil {
				httpErr = errors2.UnauthorizedException(err.Error())
			}

			c.AbortWithStatusJSON(httpErr.StatusCode, utils.NewHttpError(httpErr.Message, httpErr.Code, nil))
			return
		}

		gofortify.BlacklistTokenAndPairToken(payload)
		var user user2.User
		result := config.DB.Where("id = ?", payload.SUB).First(&user)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				err = gofortify.ErrJWTSubjectNotFound
			}
			err = errors.New("cannot find user")
		}
		c.Set("user", &user)
	}
}
