package routes

import (
	"net/http"
	"template-go-api/enums"
	middleware2 "template-go-api/middleware"

	"template-go-api/utils"

	"github.com/gin-gonic/gin"
)

func ErrorHandleWrapper(h func(*gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h(c); err != nil {
			c.Error(err)
			c.Abort()
		}
	}
}

func RegisterRoute(router *gin.Engine, c *Container) {
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, utils.NewHttpError("Route does not exist", enums.ERR_ROUTE_NOT_FOUND, nil))
	})

	api := router.Group("/api")

	{
		api.POST("/auth/authenticate", ErrorHandleWrapper(c.AuthHandler.Login))
		api.POST("/auth/refresh", middleware2.RefreshMiddleware(), ErrorHandleWrapper(c.AuthHandler.Refresh))
		api.POST("/auth/logout", middleware2.AuthMiddleware(), ErrorHandleWrapper(c.AuthHandler.Logout))
		api.GET("/auth/me", middleware2.AuthMiddleware(), ErrorHandleWrapper(c.AuthHandler.Me))
	}

	{
		api.GET("master/permissions", ErrorHandleWrapper(c.PermissionHandler.Index))
		api.GET("master/roles", ErrorHandleWrapper(c.RoleHandler.MasterIndex))
	}

	{
		api.GET("rbac/roles", ErrorHandleWrapper(c.RoleHandler.Index))
		api.POST("rbac/roles", ErrorHandleWrapper(c.RoleHandler.Store))
		api.PATCH("rbac/roles/:id", ErrorHandleWrapper(c.RoleHandler.Update))
		api.DELETE("rbac/roles/:id", ErrorHandleWrapper(c.RoleHandler.Destroy))
	}
}
