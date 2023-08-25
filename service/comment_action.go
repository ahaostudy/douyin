package service

import (
	"main/dao"
	"main/model"
)

// SendComment 发送评论
func SendComment(uid uint, vid uint, commentText string) (*model.Comment, bool) {
	cid, err := dao.InsertComment(uid, vid, commentText)
	if err != nil {
		return nil, false
	}
	comment, err := dao.GetComment(cid, uid)
	if err != nil {
		return nil, false
	}
	return comment, true
}

// DeleteComment 删除评论
func DeleteComment(commentID uint) bool {
	return dao.DeleteComment(commentID) == nil
}
