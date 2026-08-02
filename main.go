package main

import (
	"log"
	"net/http"
	"template-go-api/config"
	"template-go-api/enums"
	middleware2 "template-go-api/middleware"
	"template-go-api/routes"
	"template-go-api/utils"
	"template-go-api/validator"

	"github.com/gin-gonic/gin"
	"github.com/iqbalatma/gofortify"
)

func main() {
	config.LoadEnv()
	config.LoadLogger()
	config.ConnectDB()
	config.ConnectRDB()

	gofortify.LoadJWTConfig()
	validator.RegisterUniqueColumnValidator()
	router := gin.Default()
	router.
		Use(gin.CustomRecovery(func(c *gin.Context, _ any) {
			c.AbortWithStatusJSON(http.StatusInternalServerError, utils.NewHttpError(
				"Internal server error",
				enums.ERR_INTERNAL_SERVER_ERROR,
				nil,
			))
		})).
		Use(middleware2.LoggerMiddleware()).
		Use(middleware2.CorsMiddleware()).
		Use(middleware2.ErrorHandler())

	container := routes.NewContainer(config.DB)
	routes.RegisterRoute(router, container)

	err := router.Run(":" + config.AppConfig.AppPort)
	if err != nil {
		log.Fatal(err)
	}
}
