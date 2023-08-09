package router

import (
	"main/controller"
	"main/middleware"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())

	apiRouter := r.Group("/douyin")
	{
		apiRouter.GET("/feed", controller.Feed)
	}

	return r
}
