package middleware

import (
	"errors"
	user2 "template-go-api/app/user"
	"template-go-api/config"
	errors2 "template-go-api/errors"
	"template-go-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/iqbalatma/gofortify"
	"gorm.io/gorm"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token = c.GetHeader("Authorization")
		var atv, err = c.Cookie("access_token_verifier")
		if err != nil {
			httpErr := errors2.InvalidTokenTypeException()
			c.AbortWithStatusJSON(httpErr.StatusCode, utils.NewHttpError(httpErr.Message, httpErr.Code, nil))
			return
		}
		payload, err := gofortify.ValidateAccessToken(
			&token,
			&atv,
		)

		if err != nil {
			var httpErr *utils.HTTPError

			switch {
			case errors.Is(err, gofortify.ErrInvalidTokenType):
				httpErr = errors2.InvalidTokenTypeException()
			case errors.Is(err, gofortify.ErrExpiredToken), errors.Is(err, jwt.ErrTokenExpired):
				httpErr = errors2.UnauthorizedException("Token is expired")
			default:
				httpErr = errors2.UnauthorizedException(err.Error())
			}

			c.AbortWithStatusJSON(httpErr.StatusCode, utils.NewHttpError(httpErr.Message, httpErr.Code, nil))
			return
		}

		var user user2.User
		result := config.DB.Preload("Roles.Permissions").Where("id = ?", payload.SUB).First(&user)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				err = gofortify.ErrJWTSubjectNotFound
			}
			err = errors.New("cannot find user")
		}
		c.Set("user", &user)
		c.Set("payload", payload)
	}
}
