package dao

import (
	"gorm.io/gorm"
	"main/model"
	"time"
)

func DelCommon(commentId uint) {
	comment := model.Comment{
		ID: commentId,
	}
	DB.Delete(&comment)
}

func SendComment(uId uint, cId uint) (*model.Comment, error) {
	var comment *model.Comment
	err := DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Table("users u").Select("u.*, "+
			"COUNT(DISTINCT v.id) work_count,"+
			"COUNT(DISTINCT l.id) total_favorited,"+
			"(SELECT COUNT(*) FROM likes l WHERE l.user_id = u.id) favorite_count, "+
			"(SELECT COUNT(*) FROM follows f WHERE f.follower_id = u.id) follow_count, "+
			"(SELECT COUNT(*) FROM follows f WHERE f.user_id = u.id) follower_count, "+
			"EXISTS(SELECT * FROM follows f WHERE f.user_id = u.id AND f.follower_id = ?) is_follow", uId).
			Joins("LEFT JOIN videos v ON u.id = v.author_id").
			Joins("LEFT JOIN likes l ON v.id = l.video_id").
			Group("u.id")
	}).
		Select("comments.*").
		Where("comments.id = ?", cId).
		Find(&comment).Error

	return comment, err
}

func AddComment(uId uint, vId uint, commentText string) uint {
	comment := model.Comment{
		UserID:      uId,
		VideoID:     vId,
		CommentText: commentText,
		CreatedAt:   time.Now(),
	}
	DB.Create(&comment)
	return comment.ID
}

func GetListComment(vId uint, uId uint) ([]*model.Comment, error) {
	var commentList []*model.Comment
	err := DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Table("users u").Select("u.*, "+
			"COUNT(DISTINCT v.id) work_count,"+
			"COUNT(DISTINCT l.id) total_favorited,"+
			"(SELECT COUNT(*) FROM likes l WHERE l.user_id = u.id) favorite_count, "+
			"(SELECT COUNT(*) FROM follows f WHERE f.follower_id = u.id) follow_count, "+
			"(SELECT COUNT(*) FROM follows f WHERE f.user_id = u.id) follower_count, "+
			"EXISTS(SELECT * FROM follows f WHERE f.user_id = u.id AND f.follower_id = ?) is_follow", uId).
			Joins("LEFT JOIN videos v ON u.id = v.author_id").
			Joins("LEFT JOIN likes l ON v.id = l.video_id").
			Group("u.id")
	}).
		Select("comments.*").
		Where("comments.video_id", vId).
		Order("comments.created_at").
		Find(&commentList).Error

	return commentList, err
}
