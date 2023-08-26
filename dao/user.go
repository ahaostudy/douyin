package dao

import (
	"main/model"
)

// GetUserInfoByID 获取用户信息
func GetUserInfoByID(id uint) (*model.User, error) {
	user := new(model.User)
	err := DB.Select("u.*, "+
		"COUNT(DISTINCT v.id) work_count,"+
		"COUNT(DISTINCT l.id) total_favorited,"+
		"(SELECT COUNT(*) FROM likes l WHERE l.user_id = u.id) favorite_count, "+
		"(SELECT COUNT(*) FROM follows f WHERE f.follower_id = u.id) follow_count, "+
		"(SELECT COUNT(*) FROM follows f WHERE f.user_id = u.id) follower_count").
		Model(user).Table("users u").
		Joins("LEFT JOIN videos v ON v.author_id = u.id").
		Joins("LEFT JOIN likes l ON l.video_id = v.id").
		Where("u.id = ?", id).
		First(user).Error
	return user, err
}

// GetUserByID 通过ID获取用户信息
func GetUserByID(id, curID uint) (*model.User, error) {
	user := new(model.User)
	// 查询用户信息
	err := DB.Select("u.*, "+
		"COUNT(DISTINCT v.id) work_count,"+
		"COUNT(DISTINCT l.id) total_favorited,"+
		"(SELECT COUNT(*) FROM likes l WHERE l.user_id = u.id) favorite_count, "+
		"(SELECT COUNT(*) FROM follows f WHERE f.follower_id = u.id) follow_count, "+
		"(SELECT COUNT(*) FROM follows f WHERE f.user_id = u.id) follower_count, "+
		"EXISTS(SELECT * FROM follows f WHERE f.user_id = u.id AND f.follower_id = ?) is_follow", curID).
		Model(user).Table("users u").
		Joins("LEFT JOIN videos v ON v.author_id = u.id").
		Joins("LEFT JOIN likes l ON l.video_id = v.id").
		Where("u.id = ?", id).
		First(user).Error
	return user, err
}

// GetUserByUsername 通过用户名获取用户
func GetUserByUsername(username string) (*model.User, error) {
	user := new(model.User)
	err := DB.Where("username = ?", username).First(user).Error
	return user, err
}

// InsertUser 插入一条用户信息
func InsertUser(user *model.User) (*model.User, error) {
	err := DB.Create(&user).Error
	return user, err
}

// GetBasicFollowList 获取用户关注ID列表，不进行联表获取用户详细信息
func GetBasicFollowList(id uint) ([]*model.Follow, error) {
	var followList []*model.Follow
	err := DB.Where("follower_id = ?", id).Find(&followList).Error
	return followList, err
}

// GetFollowList 获取用户的关注列表
func GetFollowList(id uint, curID uint) ([]*model.User, error) {
	var followList []*model.User
	err := DB.Select("u.*, "+
		"COUNT(DISTINCT v.id) work_count,"+
		"COUNT(DISTINCT l.id) total_favorited,"+
		"(SELECT COUNT(*) FROM likes l WHERE l.user_id = u.id) favorite_count, "+
		"(SELECT COUNT(*) FROM follows f WHERE f.follower_id = u.id) follow_count, "+
		"(SELECT COUNT(*) FROM follows f WHERE f.user_id = u.id) follower_count, "+
		"EXISTS(SELECT * FROM follows f WHERE f.follower_id = ? AND f.user_id = u.id) is_follow", curID).
		Table("users u").
		Joins("LEFT JOIN videos v ON v.author_id = u.id").
		Joins("LEFT JOIN likes l ON l.video_id = v.id").
		Joins("LEFT JOIN follows f ON f.user_id = u.id").
		Where("f.follower_id = ?", id).
		Group("u.id").
		Find(&followList).Error
	return followList, err
}

// GetBasicFollowerList 获取用户粉丝ID列表，不进行联表获取用户详细信息
func GetBasicFollowerList(id uint) ([]*model.Follow, error) {
	var followerList []*model.Follow
	err := DB.Where("user_id = ?", id).Find(&followerList).Error
	return followerList, err
}

// GetFollowerList 获取用户的粉丝列表
func GetFollowerList(id uint, curID uint) ([]*model.User, error) {
	var followerList []*model.User
	err := DB.Select("u.*, "+
		"COUNT(DISTINCT v.id) work_count,"+
		"COUNT(DISTINCT l.id) total_favorited,"+
		"(SELECT COUNT(*) FROM likes l WHERE l.user_id = u.id) favorite_count, "+
		"(SELECT COUNT(*) FROM follows f WHERE f.follower_id = u.id) follow_count, "+
		"(SELECT COUNT(*) FROM follows f WHERE f.user_id = u.id) follower_count, "+
		"EXISTS(SELECT * FROM follows f WHERE f.follower_id = u.id AND f.user_id = ?) is_follow", curID).
		Table("users u").
		Joins("LEFT JOIN videos v ON v.author_id = u.id").
		Joins("LEFT JOIN likes l ON l.video_id = v.id").
		Joins("LEFT JOIN follows f ON f.follower_id = u.id").
		Where("f.user_id = ?", id).
		Group("u.id").
		Find(&followerList).Error
	return followerList, err
}

// GetFriendList 获取用户的朋友列表
func GetFriendList(id uint, curID uint) ([]*model.User, error) {
	var followerList []*model.User
	err := DB.Select("u.*, "+
		"COUNT(DISTINCT v.id) work_count,"+
		"COUNT(DISTINCT l.id) total_favorited,"+
		"(SELECT COUNT(*) FROM likes l WHERE l.user_id = u.id) favorite_count, "+
		"(SELECT COUNT(*) FROM follows f WHERE f.follower_id = u.id) follow_count, "+
		"(SELECT COUNT(*) FROM follows f WHERE f.user_id = u.id) follower_count, "+
		"EXISTS(SELECT * FROM follows f WHERE f.follower_id = u.id AND f.user_id = ?) is_follow", curID).
		Table("users u").
		Joins("LEFT JOIN videos v ON v.author_id = u.id").
		Joins("LEFT JOIN likes l ON l.video_id = v.id").
		Joins("JOIN follows f").
		Where("EXISTS(SELECT * FROM follows f WHERE f.user_id = u.id AND f.follower_id = ?) AND "+
			"EXISTS(SELECT * FROM follows f WHERE f.user_id = ? AND f.follower_id = u.id) AND "+
			"f.user_id = ?", id, id, id).
		Group("u.id").
		Find(&followerList).Error
	return followerList, err
}
