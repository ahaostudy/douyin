package service

import (
	"main/dao"
	"main/model"
	"time"
)

// SavaFile 保存上传的视频数据到数据库
func SavaFile(id uint, fileName string, title string) error {
	video := model.Video{
		AuthorID: id,
		Title:    title,
		PlayUrl:  fileName,
		// TODO 截取封面图
		CoverUrl:  "1.jpg",
		CreatedAt: time.Now(),
	}
	err := dao.InsertVideo(&video)
	return err
}
