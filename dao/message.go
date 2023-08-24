package dao

import (
	"fmt"
	"main/model"
	"time"
)

// InsertMessage 插入一条消息记录
func InsertMessage(message *model.Message) (*model.Message, error) {
	err := DB.Create(message).Error
	fmt.Println(message.CreatedAt.UnixMilli())
	return message, err
}

// GetMessageList 获取消息列表
func GetMessageList(userID, toUserID uint, preMsgTime time.Time) ([]*model.Message, error) {
	var messageList []*model.Message
	err := DB.Where("(from_user_id = ? AND to_user_id = ? OR to_user_id = ? AND from_user_id = ?) "+
		"AND created_at > ?", userID, toUserID, userID, toUserID, preMsgTime).Order("created_at").Find(&messageList).Error
	return messageList, err
}

// GetLatestMessageTime 获取最新的消息时间
func GetLatestMessageTime(userID, toUserID uint) (time.Time, error) {
	message := new(model.Message)
	err := DB.Where("from_user_id = ? AND to_user_id = ? OR to_user_id = ? AND from_user_id = ?",
		userID, toUserID, userID, toUserID).Order("created_at DESC").First(message).Error
	return message.CreatedAt, err
}
