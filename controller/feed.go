package controller

import (
	"main/model"
	"main/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response
	NextTime  int64          `json:"next_time,omitempty"`
	VideoList []*model.Video `json:"video_list"`
}

func Feed(c *gin.Context) {
	// 解析时间戳获取时间
	latestTimeStr := c.Query("latest_time")
	latestTime := time.Now()
	if len(latestTimeStr) > 0 {
		timeStamp, _ := strconv.ParseInt(latestTimeStr, 10, 64)
		latestTime = time.UnixMilli(timeStamp)
	}
	userID := c.GetUint("user_id")

	// 获取视频列表
	// maxCount: 30，单次获取最大视频数量为30，接口文件中要求的
	// user_id 为token上的user_id，用于获取该用户对视频的点赞和评论数据
	videoList, ok := service.GetVideoList(latestTime, 30, userID)
	if !ok {
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Server failed"},
		})
		return
	}

	// 获取最早发布视频的时间戳
	var nextTime int64
	if len(videoList) > 0 {
		nextTime = videoList[len(videoList)-1].CreatedAt.UnixMilli()
	}

	// success
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0, StatusMsg: "OK"},
		NextTime:  nextTime,
		VideoList: videoList,
	})
}
