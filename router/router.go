package router

import (
	"main/controller"
	"main/middleware/cors"
	"main/middleware/jwt"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Cors())

	apiRouter := r.Group("/douyin")
	{
		// 不需要鉴权的基本功能
		apiRouter.GET("/feed/", controller.Feed)
		apiRouter.POST("/user/register/", controller.Register)
		apiRouter.POST("/user/login/", controller.Login)

		// 使用鉴权中间件
		apiRouter.Use(jwt.Auth())

		// 需要鉴权的路由
		apiRouter.GET("/user/", controller.UserInfo)
	}

	return r
}
