package controller

import (
	"fmt"
	"main/config"
	"main/middleware/ffmpeg"
	"main/service"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

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

	// 有待优化，比如常量需不需要统一起来管理
	// 这里的文件后缀从 视频格式后缀 -> .jpg，文件名保持一致
	ext = ".jpg"
	coverfinalName := strings.Split(finalName, ".")[0] + ext
	coverSaveFile := filepath.Join(config.StaticPath, "cover", strconv.Itoa(int(uid)), coverfinalName)

	if ffmpeg.ExtractThumbnail(saveFile, coverSaveFile) != nil {
		c.JSON(http.StatusOK, PublishActionResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Submission failed"},
		})
		return
	}

	// 保存到数据库
	if service.SaveFile(uid, finalName, title) != nil {
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
