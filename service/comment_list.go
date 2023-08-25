package service

import (
	"main/dao"
	"main/model"
)

// GetCommentList 获取评论列表
func GetCommentList(vid uint, uid uint) ([]*model.Comment, bool) {
	commentList, err := dao.GetCommentList(vid, uid)
	if err != nil {
		return nil, false
	}
	return commentList, true
}
