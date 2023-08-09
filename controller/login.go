package controller

import (
	"github.com/gin-gonic/gin"
	"main/service"
	"main/utils"
	"net/http"
)

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	// 登录
	user, ok := service.Login(username, password)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "The username or password is incorrect",
		})
		return
	}

	// 生成token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "Server failed",
		})
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
