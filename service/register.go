package service

import (
	"gorm.io/gorm"
	"main/dao"
	"main/model"
	"main/utils"
)

// IsExistUser 判断用户是否已经存在
func IsExistUser(username string) bool {
	user, err := dao.GetUserByUsername(username)
	if err == gorm.ErrRecordNotFound || user == nil {
		return false
	}
	return true
}

// Register 注册用户
func Register(username, password string) (*model.User, bool) {
	// 将密码加密并添加到数据库
	if user, err := dao.InsertUser(&model.User{
		Username: username, Password: utils.MD5(password),
	}); err != nil {
		return nil, false
	} else {
		return user, true
	}
}
