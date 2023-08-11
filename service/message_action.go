package service

import "main/dao"

// InsertMessage 插入一条新消息
func InsertMessage(fromID, toID uint, content string) bool {
	return dao.InsertMessage(fromID, toID, content) == nil
}
