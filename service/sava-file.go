package service

import (
	"main/dao"
	"main/model"
	"time"
)

func SavaFile(id uint, fileName string, title string) error {
	video := model.Video{
		AuthorID:  id,
		Title:     title,
		PlayUrl:   fileName,
		CoverUrl:  "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
		CreatedAt: time.Now(),
	}
	err := dao.SaveFileOne(&video)
	return err
}
