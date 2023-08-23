package controller

import (
	"fmt"
	"main/config"
	"main/service"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PublishActionResponse struct {
	Response
}

func PublishAction(c *gin.Context) {
	uid := c.GetUint("user_id")
	title := c.PostForm("title")
	file, _ := c.FormFile("data")

	// 生成文件路径
	ext := filepath.Ext(file.Filename)

	finalName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	saveFile := filepath.Join(config.StaticPath, "play", strconv.Itoa(int(uid)), finalName)

	// 保存上传的文件到本地
	if c.SaveUploadedFile(file, saveFile) != nil {
		c.JSON(http.StatusOK, PublishActionResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Submission failed"},
		})
		return
	}

	// 保存到数据库
	if service.SavaFile(uid, finalName, title) != nil {
		c.JSON(http.StatusOK, PublishActionResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Submission failed"},
		})
		return
	}

	// success
	c.JSON(http.StatusOK, PublishActionResponse{
		Response: Response{StatusCode: 0, StatusMsg: "OK"},
	})
}
