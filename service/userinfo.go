package service

import (
	"main/dao"
	"main/model"
)

// GetUserByID 通过ID获取用户
// id 目标用户ID
// tid 当前登录的用户ID
func GetUserByID(id, curID uint) (*model.User, bool) {
	user, err := dao.GetUserByID(id, curID)
	if err != nil || user == nil {
		return nil, false
	}
	return user, true
}
