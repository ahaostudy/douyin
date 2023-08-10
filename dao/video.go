package dao

import (
	"fmt"
	"main/config"
	"main/model"
	"time"
)

// GetVideoList 获取视频列表
func GetVideoList(latestTime time.Time, maxCount int) ([]*model.Video, error) {
	var videoList []*model.Video
	err := DB.Select(
		"*",
		fmt.Sprintf("CONCAT('%s', videos.play_url) play_url", config.StaticDir),
		fmt.Sprintf("CONCAT('%s', videos.cover_url) cover_url", config.StaticDir),
	).Where("created_at <= ?", latestTime).Order("created_at").Limit(maxCount).Find(&videoList).Error
	return videoList, err
}
