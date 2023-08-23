package service

import (
	"main/dao"
	"main/model"
	"path"
	"strconv"
	"strings"
	"time"
)

// SaveFile 保存上传的视频数据到数据库
func SaveFile(id uint, fileName string, title string) error {
	coverName := strings.Split(fileName, ".")[0] + ".jpg"
	video := model.Video{
		AuthorID: id,
		Title:    title,
		PlayUrl:  path.Join("play", strconv.Itoa(int(id)), fileName),
		// TODO 截取封面图
		// 保存路径为 cover/{uid}/{cover_image}
		CoverUrl:  path.Join("cover", strconv.Itoa(int(id)), coverName),
		CreatedAt: time.Now(),
	}
	err := dao.InsertVideo(&video)
	return err
}
