package controller

import (
	"fmt"
	"main/service"
	"main/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserinfoResponse struct {
	Response
	ExclusivePwdUser utils.ExclusivePwdUser `json:"user"`
}

func Userinfo(c *gin.Context) {
	// ParseUint的第二个参数是进制，十进制；第三个参数是结果的位长度
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

	exclusived := utils.GetSaftUser(user)
	fmt.Println("excluded password ", exclusived)

	// success
	c.JSON(http.StatusOK, UserinfoResponse{
		Response:         Response{StatusCode: 0, StatusMsg: "OK"},
		ExclusivePwdUser: exclusived,
	})
}
