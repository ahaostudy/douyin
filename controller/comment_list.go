package controller

import (
	"github.com/gin-gonic/gin"
	"main/model"
	"main/service"
	"net/http"
	"strconv"
)

type CommentListResponse struct {
	Response
	CommentList []*model.Comment `json:"comment_list"`
}

// CommentList 获取评论列表
func CommentList(c *gin.Context) {
	// 解析参数
	videoID, _ := strconv.Atoi(c.Query("video_id"))
	vid := uint(videoID)
	uid := c.GetUint("user_id")

	// 获取评论列表
	commentList, ok := service.GetCommentList(vid, uid)
	if !ok {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Get comment list failed"},
		})
		return
	}

	// success
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0, StatusMsg: "OK"},
		CommentList: commentList,
	})

}
