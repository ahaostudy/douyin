package service

import (
	"main/dao"
	"main/model"
)

func GetFollowList(id, curID uint) ([]*model.User, bool) {
	followList, err := dao.GetFollowList(id, curID)
	return followList, err == nil
}
