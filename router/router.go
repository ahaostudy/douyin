package router

import (
	"main/controller"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
  r := gin.Default()
  
  apiRouter := r.Group("/douyin")
  {
    apiRouter.GET("/feed", controller.Feed)
  }
  
  return r
}