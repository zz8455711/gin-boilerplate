package routers

import (
	"gin-boilerplate/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes add all routing list here automatically get main router
func RegisterRoutes(route *gin.Engine) {
	route.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Route Not Found"})
	})
	route.GET("/health", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"live": "ok"}) })
	route.GET("/v1/address", controllers.GetData)
	route.GET("/v1/address/:address", controllers.GetOneData)
	route.POST("/v1/address", controllers.Create)

	//Add All route
	//TestRoutes(route)
}
