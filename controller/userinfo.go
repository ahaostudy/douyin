package controller

import (
	"github.com/gin-gonic/gin"
	"main/service"
	"net/http"
	"strconv"
)

func UserInfo(c *gin.Context) {
	qid, _ := strconv.ParseUint(c.Query("user_id"), 10, 32)
	queryUserID := uint(qid)
	tid, _ := c.Get("user_id")
	tokenUserID := tid.(uint)

	// 参数内容不一致
	if queryUserID != tokenUserID {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "Identity verification failed",
			"user":        nil,
		})
		return
	}

	// 获取用户信息
	user, ok := service.GetUserByID(queryUserID)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "The user does not exist",
			"user":        nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "OK",
		"user":        user,
	})
}
