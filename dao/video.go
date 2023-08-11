package dao

import (
	"gorm.io/gorm"
	"main/model"
	"time"
)

// GetVideoList 获取视频列表
func GetVideoList(latestTime time.Time, maxCount int, userID uint) ([]*model.Video, error) {
	var videoList []*model.Video

	err := DB.Preload("Author", func(db *gorm.DB) *gorm.DB {
		return db.Table("users u").Select("u.*, " +
			"COUNT(DISTINCT v.id) work_count," +
			"COUNT(DISTINCT lv.id) total_favorited," +
			"COUNT(DISTINCT lu.id) favorite_count").
			Joins("LEFT JOIN videos v ON u.id = v.author_id").
			Joins("LEFT JOIN likes lv ON v.id = lv.video_id").
			Joins("LEFT JOIN likes lu ON u.id = lu.id").
			Group("u.id")
	}).
		Select("videos.*, "+
			"(SELECT COUNT(*) FROM likes l WHERE videos.id = l.video_id) AS favorite_count, "+
			"(SELECT COUNT(*) FROM comments c WHERE videos.id = c.video_id) AS comment_count, "+
			"EXISTS(SELECT * FROM likes l WHERE videos.id = l.video_id AND l.user_id = ?) AS is_favorite",
			userID).
		Where("videos.created_at <= ?", latestTime).
		Order("videos.created_at").
		Limit(maxCount).
		Find(&videoList).Error

	return videoList, err
}

func GetVideoListById(uId uint, qId uint) ([]*model.Video, error) {
	var videoList []*model.Video

	err := DB.Preload("Author", func(db *gorm.DB) *gorm.DB {
		return db.Table("users u").Select("u.*, " +
			"COUNT(DISTINCT v.id) work_count," +
			"COUNT(DISTINCT lv.id) total_favorited," +
			"COUNT(DISTINCT lu.id) favorite_count").
			Joins("LEFT JOIN videos v ON u.id = v.author_id").
			Joins("LEFT JOIN likes lv ON v.id = lv.video_id").
			Joins("LEFT JOIN likes lu ON u.id = lu.id").
			Group("u.id")
	}).
		Select("videos.*, "+
			"(SELECT COUNT(*) FROM likes l WHERE videos.id = l.video_id) AS favorite_count, "+
			"(SELECT COUNT(*) FROM comments c WHERE videos.id = c.video_id) AS comment_count, "+
			"EXISTS(SELECT * FROM likes l WHERE videos.id = l.video_id AND l.user_id = ?) AS is_favorite",
			qId).
		Where("videos.author_id = ?", uId).
		Order("videos.created_at").
		Find(&videoList).Error

	return videoList, err
}
