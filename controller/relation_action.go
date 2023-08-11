package controller

import (
	"github.com/gin-gonic/gin"
	"main/service"
	"net/http"
	"strconv"
)

type RelationActionResponse struct {
	Response
}

func RelationAction(c *gin.Context) {
	// 解析请求参数
	userID := c.GetUint("user_id")
	tui, _ := strconv.ParseUint(c.Query("to_user_id"), 10, 32)
	toUserID := uint(tui)
	at, _ := strconv.ParseInt(c.Query("action_type"), 10, 32)
	actionType := int(at)

	// 关注业务
	if !service.RelationAction(userID, toUserID, actionType) {
		c.JSON(http.StatusOK, RelationActionResponse{
			Response: Response{StatusCode: 1, StatusMsg: "relation action failed"},
		})
		return
	}

	// success
	c.JSON(http.StatusOK, RelationActionResponse{
		Response: Response{StatusCode: 0, StatusMsg: "OK"},
	})
}
