package controller

import (
	"github.com/gin-gonic/gin"
	"main/model"
	"main/service"
	"net/http"
	"strconv"
)

type FollowListResponse struct {
	Response
	UserList []*model.User `json:"user_list"`
}

func FollowList(c *gin.Context) {
	qid, _ := strconv.ParseUint(c.Query("user_id"), 10, 32)
	queryUserID := uint(qid)
	tokenUserID := c.GetUint("user_id")

	followList, ok := service.GetFollowList(queryUserID, tokenUserID)
	if !ok {
		c.JSON(http.StatusOK, FollowListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Get follow list failed"},
		})
		return
	}

	c.JSON(http.StatusOK, FollowListResponse{
		Response: Response{StatusCode: 0, StatusMsg: "OK"},
		UserList: followList,
	})
}
