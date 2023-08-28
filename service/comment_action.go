package service

import (
	"main/dao"
	"main/model"
)

// SendComment 发送评论
func SendComment(uid uint, vid uint, commentText string) (*model.Comment, bool) {
	// 先将当前视频的所有评论都加载到redis中
	if _, flag := GetCommentList(vid, uid);flag == false{
		return nil, false
	}

	if()

	// 插入一条评论记录
	cid, err := dao.InsertComment(uid, vid, commentText)
	if err != nil {
		return nil, false
	}

	// 新数据要先判断在redis中有没有对应的视频vid的记录
	// 直接将数据库最新的这个视频vid的所有评论都查询出来重新加入redis中
	// 两种情况：1. redis原本没有这个key；2. redis有这个key，那就更新
	// 更新redis中对应的视频 评论列表
	GetCommentList(vid, uid)

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
