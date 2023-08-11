package dao

import (
	"main/model"
)

// InsertMessage 插入一条消息记录
func InsertMessage(fromID, toID uint, content string) error {
	return DB.Create(&model.Message{
		FromUserID: fromID,
		ToUserID:   toID,
		Content:    content,
	}).Error
}

// GetMessageList 获取消息列表
func GetMessageList(fromID, toID uint) ([]*model.Message, error) {
	var messageList []*model.Message
	err := DB.Where("from_user_id = ? AND to_user_id = ?", fromID, toID).Find(&messageList).Error
	return messageList, err
}
