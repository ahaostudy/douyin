package controller

import (
	"github.com/gin-gonic/gin"
	"main/model"
	"main/service"
	"net/http"
	"strconv"
)

type MessageChatResponse struct {
	Response
	MessageList []*model.Message `json:"message_list"`
}

func MessageChat(c *gin.Context) {
	// 解析请求参数
	userID := c.GetUint("user_id")
	tui, _ := strconv.ParseUint(c.Query("to_user_id"), 10, 32)
	toUserID := uint(tui)

	// 获取聊天记录
	messageList, ok := service.GetMessageList(userID, toUserID)
	if !ok {
		c.JSON(http.StatusOK, MessageChatResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Get message list failed"},
		})
		return
	}

	// success
	c.JSON(http.StatusOK, MessageChatResponse{
		Response:    Response{StatusCode: 0, StatusMsg: "OK"},
		MessageList: messageList,
	})
}
