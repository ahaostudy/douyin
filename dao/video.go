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
		return db.Table("users u").Select("u.*, "+
			"COUNT(DISTINCT v.id) work_count,"+
			"COUNT(DISTINCT lv.id) total_favorited,"+
			"(SELECT COUNT(*) FROM likes l WHERE l.user_id = u.id) favorite_count, "+
			"(SELECT COUNT(*) FROM follows f WHERE f.follower_id = u.id) follow_count, "+
			"(SELECT COUNT(*) FROM follows f WHERE f.user_id = u.id) follower_count, "+
			"EXISTS(SELECT * FROM follows f WHERE f.user_id = u.id AND f.follower_id = ?) is_follow", userID).
			Joins("LEFT JOIN videos v ON u.id = v.author_id").
			Joins("LEFT JOIN likes lv ON v.id = lv.video_id").
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

// GetVideoListByAuthorID 获取作品列表
func GetVideoListByAuthorID(authorID uint, curID uint) ([]*model.Video, error) {
	var videoList []*model.Video

	err := DB.Preload("Author", func(db *gorm.DB) *gorm.DB {
		return db.Table("users u").Select("u.*, "+
			"COUNT(DISTINCT v.id) work_count,"+
			"COUNT(DISTINCT lv.id) total_favorited,"+
			"(SELECT COUNT(*) FROM likes l WHERE l.user_id = u.id) favorite_count, "+
			"(SELECT COUNT(*) FROM follows f WHERE f.follower_id = u.id) follow_count, "+
			"(SELECT COUNT(*) FROM follows f WHERE f.user_id = u.id) follower_count, "+
			"EXISTS(SELECT * FROM follows f WHERE f.user_id = u.id AND f.follower_id = ?) is_follow", curID).
			Joins("LEFT JOIN videos v ON u.id = v.authorID").
			Joins("LEFT JOIN likes lv ON v.id = lv.video_id").
			Group("u.id")
	}).
		Select("videos.*, "+
			"(SELECT COUNT(*) FROM likes l WHERE videos.id = l.video_id) AS favorite_count, "+
			"(SELECT COUNT(*) FROM comments c WHERE videos.id = c.video_id) AS comment_count, "+
			"EXISTS(SELECT * FROM likes l WHERE videos.id = l.video_id AND l.user_id = ?) AS is_favorite",
			authorID).
		Where("videos.authorID = ?", authorID).
		Order("videos.created_at").
		Find(&videoList).Error

	return videoList, err
}

// GetVideoListByLike 获取喜欢的视频列表
func GetVideoListByLike(id, curID uint) ([]*model.Video, error) {
	var videoList []*model.Video

	err := DB.Preload("Author", func(db *gorm.DB) *gorm.DB {
		return db.Table("users u").Select("u.*, "+
			"COUNT(DISTINCT v.id) work_count,"+
			"COUNT(DISTINCT lv.id) total_favorited,"+
			"(SELECT COUNT(*) FROM likes l WHERE l.user_id = u.id) favorite_count, "+
			"(SELECT COUNT(*) FROM follows f WHERE f.follower_id = u.id) follow_count, "+
			"(SELECT COUNT(*) FROM follows f WHERE f.user_id = u.id) follower_count, "+
			"EXISTS(SELECT * FROM follows f WHERE f.user_id = u.id AND f.follower_id = ?) is_follow", curID).
			Joins("LEFT JOIN videos v ON u.id = v.author_id").
			Joins("LEFT JOIN likes lv ON v.id = lv.video_id").
			Group("u.id")
	}).
		Select("videos.*, "+
			"(SELECT COUNT(*) FROM likes l WHERE videos.id = l.video_id) AS favorite_count, "+
			"(SELECT COUNT(*) FROM comments c WHERE videos.id = c.video_id) AS comment_count, "+
			"true is_favorite").
		Where("EXISTS(SELECT * FROM likes l WHERE videos.id = l.video_id AND l.user_id = ?)", id).
		Order("videos.created_at").
		Find(&videoList).Error

	return videoList, err
}

// InsertVideo 插入一条视频数据
func InsertVideo(video *model.Video) error {
	err := DB.Create(video).Error
	return err
}
