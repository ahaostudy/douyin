package service

import (
	"main/dao"
	"main/model"
	"time"
)

// GetVideoList 获取视频列表
func GetVideoList(latestTime time.Time, maxCount int) ([]*model.Video, bool) {
	videoList, err := dao.GetVideoList(latestTime, maxCount)
	if err != nil {
		return nil, false
	}
	return videoList, true
}
