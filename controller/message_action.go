package controller

import (
	"github.com/gin-gonic/gin"
	"main/service"
	"net/http"
	"strconv"
)

type MessageActionResponse struct {
	Response
}

func MessageAction(c *gin.Context) {
	// 解析请求参数
	userID := c.GetUint("user_id")
	tui, _ := strconv.ParseUint(c.Query("to_user_id"), 10, 32)
	toUserID := uint(tui)
	Content := c.Query("content")

	// 检测操作类型参数是否正确
	if c.Query("action_type") != "1" {
		c.JSON(http.StatusOK, MessageActionResponse{
			Response: Response{StatusCode: 1, StatusMsg: "The parameter is incorrect"},
		})
		return
	}

	// 发送消息
	if !service.SendMessage(userID, toUserID, Content) {
		c.JSON(http.StatusOK, MessageActionResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Message action failed"},
		})
		return
	}

	// success
	c.JSON(http.StatusOK, MessageActionResponse{
		Response: Response{StatusCode: 0, StatusMsg: "OK"},
	})
}
