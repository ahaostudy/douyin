package controller

import (
	"github.com/gin-gonic/gin"
	"main/model"
	"main/service"
	"strconv"
)

type ResponseComment struct {
	Response
	Comment *model.Comment `json:"comment"`
}

type ListComment struct {
	Response
	CommentList []*model.Comment `json:"comment_list"`
}

func SendComment(c *gin.Context) {
	uId := c.GetUint("user_id")
	videoId, _ := strconv.Atoi(c.Query("video_id"))
	vId := uint(videoId)
	actionType := c.Query("action_type")
	if actionType == "1" {
		commentText := c.Query("comment_text")
		comment, ok := service.SendComment(uId, vId, commentText)
		if !ok {
			c.JSON(200, Response{
				StatusCode: 1,
				StatusMsg:  "send failed",
			})
		}
		c.JSON(200, ResponseComment{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "send success",
			},
			Comment: comment,
		})
	} else {
		comment := c.Query("comment_id")
		commentId, _ := strconv.Atoi(comment)
		ok := service.DelComment(uint(commentId))
		if !ok {
			c.JSON(200, Response{
				StatusCode: 1,
				StatusMsg:  "delete failed",
			})
		}
		c.JSON(200, ResponseComment{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "delete success",
			},
		})
	}
}

func GetListComment(c *gin.Context) {
	videoId, _ := strconv.Atoi(c.Query("video_id"))
	vId := uint(videoId)
	uId := c.GetUint("user_id")
	listComment, ok := service.GetListComment(vId, uId)
	if !ok {
		c.JSON(200, ListComment{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "getList failed",
			},
			CommentList: nil,
		})
		return
	}
	c.JSON(200, ListComment{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "getList success",
		},
		CommentList: listComment,
	})

}
