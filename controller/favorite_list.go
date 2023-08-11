package controller

import (
	"github.com/gin-gonic/gin"
	"main/model"
	"main/service"
	"net/http"
	"strconv"
)

type FavoriteListResponse struct {
	Response
	VideoList []*model.Video `json:"video_list"`
}

// FavoriteList 喜欢列表
func FavoriteList(c *gin.Context) {
	qid, _ := strconv.ParseUint(c.Query("user_id"), 10, 32)
	queryUserID := uint(qid)
	tokenUserID := c.GetUint("user_id")

	// ID不一致
	if queryUserID != tokenUserID {
		c.JSON(http.StatusOK, FavoriteListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Identity verification failed"},
		})
		return
	}

	// 获取喜欢的视频列表
	videoList, ok := service.GetFavoriteList(queryUserID)
	if !ok {
		c.JSON(http.StatusOK, FavoriteListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Server failed"},
		})
		return
	}

	// success
	c.JSON(http.StatusOK, FavoriteListResponse{
		Response:  Response{StatusCode: 1, StatusMsg: "OK"},
		VideoList: videoList,
	})
}
