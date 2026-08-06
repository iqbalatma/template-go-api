package main

import (
	"log"
	"net/http"
	"template-go-api/app/media"
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

	media.Boot(media.Config{
		DefaultDisk:  config.AppConfig.MediaDisk,
		LocalRoot:    config.AppConfig.MediaRoot,
		LocalBaseURL: config.AppConfig.MediaURLPrefix,
	})

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

	// File pada disk lokal disajikan langsung oleh aplikasi. Bila nanti pindah
	// ke S3, baris ini bisa dilepas karena URL-nya datang dari disk tersebut.
	router.Static(config.AppConfig.MediaURLPrefix, config.AppConfig.MediaRoot)

	container := routes.NewContainer(config.DB)
	routes.RegisterRoute(router, container)

	err := router.Run(":" + config.AppConfig.AppPort)
	if err != nil {
		log.Fatal(err)
	}
}
