package service

import (
	"main/dao"
	"main/model"
)

// GetUserByID 通过ID获取用户
func GetUserByID(id uint) (*model.User, bool) {
	user, err := dao.GetUserByID(id)
	if err != nil || user == nil {
		return nil, false
	}
	return user, true
}
