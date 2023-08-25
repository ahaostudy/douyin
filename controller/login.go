package controller

import (
	"github.com/gin-gonic/gin"
	"main/service"
	"main/utils"
	"net/http"
)

type LoginResponse struct {
	Response
	UserID uint   `json:"user_id,omitempty"`
	Token  string `json:"token,omitempty"`
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	// 登录验证
	user, ok := service.Login(username, password)
	if !ok {
		c.JSON(http.StatusOK, LoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "The username or password is incorrect"},
		})
		return
	}

	// 生成token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusOK, LoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Server failed"},
		})
		return
	}

	// success
	c.JSON(http.StatusOK, LoginResponse{
		Response: Response{StatusCode: 0, StatusMsg: "OK"},
		UserID:   user.ID,
		Token:    token,
	})
}
