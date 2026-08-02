package middleware

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"template-go-api/enums"
	"template-go-api/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		var httpResponse *utils.HTTPError

		if ginErr := c.Errors.Last(); ginErr != nil {
			originalErr := ginErr.Err

			logError(c, originalErr)
			fmt.Printf(utils.ANSI_RED+"Error : %s. Error type: %T\n", originalErr, originalErr)

			if errors.Is(originalErr, io.EOF) {
				fmt.Println(originalErr)
				c.AbortWithStatusJSON(http.StatusBadRequest, utils.NewHttpError(
					"Empty request body",
					enums.ERR_BAD_REQUEST,
					nil,
				))
				return
			}

			if errors.Is(originalErr, gorm.ErrRecordNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, utils.NewHttpError("Data not found", enums.ERR_NOT_FOUND, nil))
				return
			}

			if errors.As(originalErr, &httpResponse) {
				c.AbortWithStatusJSON(httpResponse.StatusCode, utils.NewHttpError(httpResponse.Message, httpResponse.Code, nil))
				return
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, utils.NewHttpError(originalErr.Error(), enums.ERR_INTERNAL_SERVER_ERROR, nil))
		}
	}
}

func logError(c *gin.Context, err error) {
	Logger(c).Error(err)
}
