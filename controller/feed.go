package controller

import (
	"github.com/gin-gonic/gin"
	"main/model"
	"main/service"
	"net/http"
	"strconv"
	"time"
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

	// TODO 解析token
	//token := c.Query("token")

	// 获取视频列表
	// maxCount: 30，单次获取最大视频数量为30，接口文件中要求的
	videoList, ok := service.GetVideoList(latestTime, 30)
	if !ok {
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Server failed"},
		})
		return
	}

	var nextTime int64
	if len(videoList) > 0 {
		nextTime = videoList[0].CreatedAt.UnixMilli()
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0, StatusMsg: "OK"},
		NextTime:  nextTime,
		VideoList: videoList,
	})
}
