package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"main/service"
	"main/utils"
	"net/http"
)

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	// 判断是否已经注册
	if service.IsExistUser(username) {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "Username already exists",
		})
		return
	}

	// 注册用户
	user, ok := service.Register(username, password)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "Server failed",
		})
		return
	}

	// 生成Token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "Server failed",
		})
		log.Println(err.Error())
		return
	}

	// success
	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "OK",
		"user_id":     user.ID,
		"token":       token,
	})
}
