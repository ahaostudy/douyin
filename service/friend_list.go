package service

import (
	"main/dao"
	"main/model"
)

func GetFriendList(id, curID uint) ([]*model.User, bool) {
	friendList, err := dao.GetFriendList(id, curID)
	return friendList, err == nil
}
