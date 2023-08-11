package service

import (
	"main/dao"
	"main/model"
)

func GetFollowerList(id, curID uint) ([]*model.User, bool) {
	followList, err := dao.GetFollowerList(id, curID)
	return followList, err == nil
}
