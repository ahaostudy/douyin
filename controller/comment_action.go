package controller

import (
	"github.com/gin-gonic/gin"
	"main/model"
	"main/service"
	"net/http"
	"strconv"
)

type CommentActionResponse struct {
	Response
	Comment *model.Comment `json:"comment"`
}

// CommentAction 评论操作
func CommentAction(c *gin.Context) {
	// 解析参数
	uid := c.GetUint("user_id")
	videoID, _ := strconv.Atoi(c.Query("video_id"))
	vid := uint(videoID)
	actionType := c.Query("action_type")

	if actionType == "1" {
		// 发送评论
		commentText := c.Query("comment_text")

		comment, ok := service.SendComment(uid, vid, commentText)
		if !ok {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: 1, StatusMsg: "Send comment failed"},
			})
			return
		}

		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: 0, StatusMsg: "OK"},
			Comment:  comment,
		})
	} else {
		// 删除评论
		comment := c.Query("comment_id")
		commentID, _ := strconv.Atoi(comment)

		if ok := service.DeleteComment(uint(commentID)); !ok {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: 1, StatusMsg: "Delete comment failed"},
			})
			return
		}

		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: 0, StatusMsg: "OK"},
		})
	}
}
