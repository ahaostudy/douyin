package service

import (
	"main/dao"
	"main/model"
)

// SendComment 发送评论
func SendComment(uid uint, vid uint, commentText string) (*model.Comment, bool) {
	// 插入一条评论记录
	cid, err := dao.InsertComment(uid, vid, commentText)
	if err != nil {
		return nil, false
	}
	// 获取评论的详细数据
	comment, err := dao.GetComment(cid, uid)
	if err != nil {
		return nil, false
	}
	return comment, true
}

// DeleteComment 删除评论
func DeleteComment(commentID uint, uid uint) bool {
	// 获取评论基础信息，判断发表评论的用户ID与当前用户ID是否一致
	comment, err := dao.GetCommentBasicInfo(commentID)
	if err != nil || comment.UserID != uid {
		return false
	}
	// 删除评论
	return dao.DeleteComment(commentID) == nil
}
