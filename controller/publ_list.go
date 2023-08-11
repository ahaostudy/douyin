package controller

import (
	"github.com/gin-gonic/gin"
	"main/model"
	"main/service"
	"net/http"
	"strconv"
)

type VideoResponse struct {
	Response
	VideoList []*model.Video `json:"video_list"`
}

func PubList(c *gin.Context) {
	sId := c.Query("user_id")
	pId, _ := strconv.Atoi(sId)
	uId := uint(pId)
	qId := c.GetUint("user_id")
	videoList, b := service.GetVideoOne(uId, qId)
	if !b {
		c.JSON(http.StatusOK, VideoResponse{
			Response{StatusCode: 1, StatusMsg: "Server failed"},
			nil,
		})
		return
	}
	c.JSON(http.StatusOK, VideoResponse{
		Response{StatusCode: 1, StatusMsg: "Server failed"},
		videoList,
	})
}
