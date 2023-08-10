package dao

import (
	"fmt"
	"github.com/spf13/viper"
	"main/model"
	"time"
)

// GetVideoList 获取视频列表
func GetVideoList(latestTime time.Time, maxCount int) ([]*model.Video, error) {
	var videoList []*model.Video
	err := DB.Select(
		"*",
		fmt.Sprintf("CONCAT('%s', videos.play_url) play_url", viper.GetString("server.static")),
		fmt.Sprintf("CONCAT('%s', videos.cover_url) cover_url", viper.GetString("server.static")),
	).Where("created_at <= ?", latestTime).Order("created_at").Limit(maxCount).Find(&videoList).Error
	return videoList, err
}
