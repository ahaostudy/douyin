package controller

import (
	"github.com/gin-gonic/gin"
	"main/service"
	"net/http"
	"strconv"
)

type FavoriteActionResponse struct {
	Response
}

// FavoriteAction 赞操作
func FavoriteAction(c *gin.Context) {
	// 解析请求参数
	userID := c.GetUint("user_id")
	vi, _ := strconv.ParseUint(c.Query("video_id"), 10, 32)
	videoID := uint(vi)
	at, _ := strconv.ParseInt(c.Query("action_type"), 10, 32)
	actionType := int(at)

	// 点赞业务
	if !service.FavoriteAction(userID, videoID, actionType) {
		c.JSON(http.StatusOK, FavoriteActionResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Favorite action failed"},
		})
		return
	}

	// success
	c.JSON(http.StatusOK, FavoriteActionResponse{
		Response: Response{StatusCode: 0, StatusMsg: "OK"},
	})
}
