package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"main/service"
	"main/utils"
	"net/http"
)

type RegisterResponse struct {
	Response
	UserID uint   `json:"user_id"`
	Token  string `json:"token"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	// 判断是否已经注册
	if service.IsExistUser(username) {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1, StatusMsg: "Username already exists",
		})
		return
	}

	// 注册用户
	user, ok := service.Register(username, password)
	if !ok {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1, StatusMsg: "Server failed",
		})
		return
	}

	// 生成Token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1, StatusMsg: "Server failed",
		})
		log.Println(err.Error())
		return
	}

	// success
	c.JSON(http.StatusOK, RegisterResponse{
		Response: Response{StatusCode: 0, StatusMsg: "OK"},
		UserID:   user.ID,
		Token:    token,
	})
}
