package dao

import (
	"main/model"
)

// GetLike 获取点赞记录
func GetLike(uid, vid uint) (*model.Like, error) {
	like := new(model.Like)
	err := DB.Where("user_id = ? AND video_id = ?", uid, vid).First(like).Error
	return like, err
}

// InsertLike 插入点赞记录
func InsertLike(uid, vid uint) error {
	return DB.Create(&model.Like{UserID: uid, VideoID: vid}).Error
}

// DeleteLike
// 删除点赞记录，返回操作的行数
func DeleteLike(uid, vid uint) error {
	return DB.Delete(new(model.Like), "user_id = ? AND video_id = ?", uid, vid).Error
}
