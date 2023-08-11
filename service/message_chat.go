package service

import (
	"main/dao"
	"main/model"
)

// GetMessageList 获取消息列表
func GetMessageList(fromID, toID uint) ([]*model.Message, bool) {
	messageList, err := dao.GetMessageList(fromID, toID)
	if err != nil {
		return nil, false
	}
	return messageList, true
}
