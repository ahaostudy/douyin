package service

import (
	"main/dao"
	"main/middleware/redis"
	"main/model"
	"strconv"
)

// SendComment 发送评论
func SendComment(uid uint, vid uint, commentText string) (*model.Comment, bool) {
	// 插入一条评论记录
	cid, err := dao.InsertComment(uid, vid, commentText)
	if err != nil {
		return nil, false
	}

	go func() {
		// 并发更新redis中的对应视频vid的list
		ctx, cancel := redis.WithTimeoutContextBySecond(3)
		defer cancel()
		key := redis.GenerateCommentKey(vid)
		if err := updateRedisComments(ctx, key, vid, uid); err != nil {
			return
		}
	}()

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
	vid := comment.VideoID

	// if err := rabbitmq.RMQUnLike.Publish(rabbitmq.GenerateUnLikeMQParam(uid, vid)); err != nil {
	// 	redis.RdbLike.SAdd(ctx, key, uid)
	// 	return false
	// }

	go func() {
		// 并发更新redis中的对应视频vid的list
		key := redis.GenerateCommentKey(vid)
		ctx, cancel := redis.WithTimeoutContextBySecond(3)
		defer cancel()

		if LoadCommentList(ctx, vid, uid) != nil {
			return
		}

		if err := redis.RdbComment.HDel(ctx, key, strconv.Itoa(int(commentID))).Err(); err != nil {
			return
		}
	}()

	// 删除评论
	return dao.DeleteComment(commentID) == nil
}
