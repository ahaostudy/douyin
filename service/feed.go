package service

import (
	"main/dao"
	"main/model"
	"time"
)

// GetVideoList 获取视频列表
func GetVideoList(latestTime time.Time, maxCount int, userID uint) ([]*model.Video, bool) {
	var videoList []*model.Video
	var err error
	videoList, err = dao.GetVideoList(latestTime, maxCount, userID)
	if err != nil {
		return nil, false
	}
	return videoList, true
}
