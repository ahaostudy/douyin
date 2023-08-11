package service

import (
	"main/dao"
	"main/model"
)

// GetFavoriteList 获取用户喜欢的视频列表
func GetFavoriteList(userID uint) ([]*model.Video, bool) {
	videoList, err := dao.GetVideoListByLike(userID)
	if err != nil {
		return nil, false
	}
	return videoList, true
}
