package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"main/config"
	"main/middleware/ffmpeg"
	"main/model"
	"main/service"
	"net/http"
	"path/filepath"
	"strconv"
)

type PublishActionResponse struct {
	Response
}

func PublishAction(c *gin.Context) {
	uid := c.GetUint("user_id")
	title := c.PostForm("title")
	file, _ := c.FormFile("data")

	// 生成文件路径
	fileName := uuid.New().String()
	ext := filepath.Ext(file.Filename)
	videoFinalName := fmt.Sprintf("%s%s", fileName, ext)
	videoSaveFile := filepath.Join(config.StaticPath, "play", strconv.Itoa(int(uid)), videoFinalName)

	// 保存上传的文件到本地
	if c.SaveUploadedFile(file, videoSaveFile) != nil {
		c.JSON(http.StatusOK, PublishActionResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Submission failed"},
		})
		return
	}

	// 有待优化，比如常量需不需要统一起来管理
	// 这里的文件后缀从 视频格式后缀 -> .jpg，文件名保持一致
	coverFinalName := fmt.Sprintf("%s.jpg", fileName)
	coverSaveFile := filepath.Join(config.StaticPath, "cover", strconv.Itoa(int(uid)), coverFinalName)

	// 截取封面图
	if ffmpeg.ExtractThumbnail(videoSaveFile, coverSaveFile) != nil {
		c.JSON(http.StatusOK, PublishActionResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Submission failed"},
		})
		return
	}

	// 更新数据库
	if service.PublishAction(&model.Video{
		AuthorID: uid, PlayUrl: videoFinalName, CoverUrl: coverFinalName, Title: title,
	}) != nil {
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
