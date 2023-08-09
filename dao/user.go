package dao

import "main/model"

// GetUserByID 通过用户ID获取用户
func GetUserByID(id uint) (*model.User, error) {
	user := new(model.User)
	err := DB.First(user, id).Error
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
