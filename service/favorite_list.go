package service

import (
	"main/dao"
	"main/model"
)

// GetFavoriteList 获取用户喜欢的视频列表
func GetFavoriteList(id, curID uint) ([]*model.Video, bool) {
	videoList, err := dao.GetVideoListByLike(id, curID)
	if err != nil {
		return nil, false
	}
	return videoList, true
}
