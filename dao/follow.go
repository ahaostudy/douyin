package dao

import "main/model"

// GetFollow 获取关注记录
func GetFollow(uid, tid uint) (*model.Follow, error) {
	follow := new(model.Follow)
	err := DB.Where("user_id = ? AND follower_id = ?", tid, uid).First(follow).Error
	return follow, err
}

// InsertFollow 插入记录记录
func InsertFollow(uid, tid uint) error {
	return DB.Create(&model.Follow{UserID: tid, FollowerID: uid}).Error
}

// DeleteFollow 删除关注记录
func DeleteFollow(uid, tid uint) error {
	return DB.Delete(new(model.Follow), "user_id = ? AND follower_id = ?", tid, uid).Error
}
