package service

import (
	"main/dao"
	"main/model"
	"time"
)

// GetMessageList 获取消息列表
func GetMessageList(fromID, toID uint, preMsgTime time.Time) ([]*model.Message, bool) {
	messageList, err := dao.GetMessageList(fromID, toID, preMsgTime)
	if err != nil {
		return nil, false
	}
	return messageList, true
}
