package router

import (
	"main/controller"
	"main/middleware/cors"
	"main/middleware/jwt"
	"os"
	"path"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	r := gin.Default()

	// 跨域处理
	r.Use(cors.Cors())

	// 配置静态路径
	workDir, _ := os.Getwd()
	r.Static("/static", path.Join(workDir, "public"))

	// douyin api
	apiRouter := r.Group("/douyin")
	{
		// 不需要鉴权的基本功能
		apiRouter.POST("/user/register/", controller.Register)
		apiRouter.POST("/user/login/", controller.Login)

		apiRouter.GET("/feed/", jwt.Parse(), controller.Feed)

		// 鉴权中间件
		apiRouter.Use(jwt.Auth())

		// 需要鉴权的路由
		apiRouter.GET("/user/", controller.Userinfo)
		apiRouter.POST("/publish/action/", controller.PubAction)
		apiRouter.GET("/publish/list/", controller.PubList)
		apiRouter.POST("/favorite/action/", controller.FavoriteAction)
		apiRouter.GET("/favorite/list/", controller.FavoriteList)
		apiRouter.POST("/relation/action/", controller.RelationAction)
		apiRouter.GET("/relation/follow/list/", controller.FollowList)
		apiRouter.GET("/relation/follower/list/", controller.FollowerList)
		apiRouter.GET("/relation/friend/list/", controller.FriendList)
		apiRouter.POST("/message/action/", controller.MessageAction)
		apiRouter.GET("/message/chat/", controller.MessageChat)
		apiRouter.POST("/comment/action/", controller.SendComment)
		apiRouter.GET("/comment/list/", controller.GetListComment)
	}

	return r
}
