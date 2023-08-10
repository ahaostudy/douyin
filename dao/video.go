package dao

import (
	"main/model"
	"time"
)

// GetVideoList 获取视频列表
func GetVideoList(latestTime time.Time, maxCount int, userID uint) ([]*model.Video, error) {
	var videoList []*model.Video

	err := DB.Preload("Author").
		Select("videos.*, "+
			"(SELECT COUNT(*) FROM likes l WHERE videos.id = l.video_id) AS favorite_count, "+
			"(SELECT COUNT(*) FROM comments c WHERE videos.id = c.video_id) AS comment_count, "+
			"EXISTS(SELECT * FROM likes l WHERE videos.id = l.video_id AND l.user_id = ?) AS is_favorite",
			userID).
		Joins("LEFT JOIN (SELECT u.*, COUNT(DISTINCT v.id) AS work_count, "+
			"COUNT(DISTINCT lv.id) AS total_favorited, COUNT(DISTINCT lu.id) AS favorite_count "+
			"FROM users u "+
			"LEFT JOIN videos v ON u.id = v.author_id "+
			"LEFT JOIN likes lv ON v.id = lv.video_id "+
			"LEFT JOIN likes lu ON u.id = lu.user_id "+
			"GROUP BY u.id) users ON users.id = videos.author_id").
		Where("videos.created_at <= ?", latestTime).
		Order("videos.created_at").
		Limit(maxCount).
		Find(&videoList).Error

	return videoList, err
}
