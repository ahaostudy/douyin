package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/model"
	"main/service"
	"net/http"
	"strconv"
)

type UserinfoResponse struct {
	Response
	User *model.User `json:"user"`
}

func Userinfo(c *gin.Context) {
	qid, _ := strconv.ParseUint(c.Query("user_id"), 10, 32)
	queryUserID := uint(qid)
	tokenUserID := c.GetUint("user_id")

	// 获取用户信息
	user, ok := service.GetUserByID(queryUserID, tokenUserID)
	fmt.Println(user, ok)
	if !ok {
		c.JSON(http.StatusOK, UserinfoResponse{
			Response: Response{StatusCode: 1, StatusMsg: "The user does not exist"},
		})
		return
	}

	// success
	c.JSON(http.StatusOK, UserinfoResponse{
		Response: Response{StatusCode: 0, StatusMsg: "OK"},
		User:     user,
	})
}
