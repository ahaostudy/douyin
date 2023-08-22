package service

import (
	"main/dao"
	"main/model"
)

// GetVideoListByAuthorID 获取一条视频
func GetVideoListByAuthorID(authorID uint, curID uint) ([]*model.Video, bool) {
	videoList, err := dao.GetVideoListByAuthorID(authorID, curID)
	if err != nil {
		return nil, false
	}
	return videoList, true
}
