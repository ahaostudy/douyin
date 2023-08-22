package service

import (
	"main/dao"
	"main/model"
	"path"
	"strconv"
	"time"
)

// SavaFile 保存上传的视频数据到数据库
func SavaFile(id uint, fileName string, title string) error {
	video := model.Video{
		AuthorID: id,
		Title:    title,
		PlayUrl:  path.Join("play", strconv.Itoa(int(id)), fileName),
		// TODO 截取封面图
		// 保存路径为 cover/{uid}/{cover_image}
		CoverUrl:  "cover/cover.jpg",
		CreatedAt: time.Now(),
	}
	err := dao.InsertVideo(&video)
	return err
}
