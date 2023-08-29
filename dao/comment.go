package dao

import (
	"main/model"
	"time"

	"gorm.io/gorm"
)

func DeleteComment(commentID uint) error {
	comment := model.Comment{
		ID: commentID,
	}
	return DB.Delete(&comment).Error
}

func GetCommentBasicInfo(id uint) (*model.Comment, error) {
	comment := &model.Comment{} // 创建 Comment 结构体的实例

	err := DB.First(comment, id).Error
	return comment, err
}

func GetComment(cid, uid uint) (*model.Comment, error) {
	var comment *model.Comment
	err := DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Table("users u").Select("u.*, "+
			"COUNT(DISTINCT v.id) work_count,"+
			"COUNT(DISTINCT l.id) total_favorited,"+
			"(SELECT COUNT(*) FROM likes l WHERE l.user_id = u.id) favorite_count, "+
			"(SELECT COUNT(*) FROM follows f WHERE f.follower_id = u.id) follow_count, "+
			"(SELECT COUNT(*) FROM follows f WHERE f.user_id = u.id) follower_count, "+
			"EXISTS(SELECT * FROM follows f WHERE f.user_id = u.id AND f.follower_id = ?) is_follow", uid).
			Joins("LEFT JOIN videos v ON u.id = v.author_id").
			Joins("LEFT JOIN likes l ON v.id = l.video_id").
			Group("u.id")
	}).
		Select("comments.*").
		Where("comments.id = ?", cid).
		First(&comment).Error

	return comment, err
}

func InsertComment(uid uint, vid uint, commentText string) (*model.Comment, error) {
	comment := &model.Comment{
		UserID:      uid,
		VideoID:     vid,
		CommentText: commentText,
		CreatedAt:   time.Now(),
	}
	err := DB.Create(comment).Error
	return comment, err
}

func GetCommentList(vid uint, uid uint) ([]*model.Comment, error) {
	var commentList []*model.Comment
	err := DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Table("users u").Select("u.*, "+
			"COUNT(DISTINCT v.id) work_count,"+
			"COUNT(DISTINCT l.id) total_favorited,"+
			"(SELECT COUNT(*) FROM likes l WHERE l.user_id = u.id) favorite_count, "+
			"(SELECT COUNT(*) FROM follows f WHERE f.follower_id = u.id) follow_count, "+
			"(SELECT COUNT(*) FROM follows f WHERE f.user_id = u.id) follower_count, "+
			"EXISTS(SELECT * FROM follows f WHERE f.user_id = u.id AND f.follower_id = ?) is_follow", uid).
			Joins("LEFT JOIN videos v ON u.id = v.author_id").
			Joins("LEFT JOIN likes l ON v.id = l.video_id").
			Group("u.id")
	}).
		Select("comments.*").
		Where("comments.video_id", vid).
		Order("comments.created_at DESC").
		Find(&commentList).Error

	return commentList, err
}
