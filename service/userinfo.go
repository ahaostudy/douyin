package service

import (
	"main/dao"
	"main/model"
)

func GetUserByID(id uint) (*model.User, bool) {
	user, err := dao.GetUserByID(id)
	if err != nil {
		return nil, false
	}
	return user, true
}
