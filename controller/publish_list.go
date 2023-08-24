package controller

import (
	"main/model"
	"main/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VideoResponse struct {
	Response
	VideoList []*model.Video `json:"video_list,omitempty"`
}

func PublishList(c *gin.Context) {
	qid, _ := strconv.ParseUint(c.Query("user_id"), 10, 32)
	tid := c.GetUint("user_id")

	// 获取用户的作品列表
	videoList, ok := service.GetVideoListByAuthorID(uint(qid), tid)
	if !ok {
		c.JSON(http.StatusOK, VideoResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Server failed"},
		})
		return
	}

	// success
	c.JSON(http.StatusOK, VideoResponse{
		Response:  Response{StatusCode: 0, StatusMsg: "OK"},
		VideoList: videoList,
	})
}
