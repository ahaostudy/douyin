package dao

import (
	"fmt"
	"github.com/spf13/viper"
	"main/model"
)

// GetUserByID 通过ID获取用户信息
func GetUserByID(id uint) (*model.User, error) {
	user := new(model.User)
	// 联表查询用户基本信息、作品数、获赞数、点赞数
	err := DB.Select(
		"u.*",
		fmt.Sprintf("CONCAT('%s', u.avatar) avatar", viper.GetString("server.static")),
		fmt.Sprintf("CONCAT('%s', u.background_image) background_image", viper.GetString("server.static")),
		"COUNT(DISTINCT v.id) work_count",
		"COUNT(DISTINCT lv.id) total_favorited",
		"COUNT(DISTINCT lu.id) favorite_count",
	).Model(user).Table("users u").
		Joins("LEFT JOIN videos v ON u.id = v.author_id").
		Joins("LEFT JOIN likes lv ON v.id = lv.video_id").
		Joins("LEFT JOIN likes lu ON u.id = lu.id").
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
