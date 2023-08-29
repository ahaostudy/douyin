package service

import (
	"encoding/json"
	"main/dao"
	"main/middleware/rabbitmq"
	"main/middleware/redis"
	"main/model"
	"strconv"
)

// SendComment 发送评论
func SendComment(uid uint, vid uint, commentText string) (*model.Comment, bool) {
	// 插入一条评论记录
	comment, err := dao.InsertComment(uid, vid, commentText)
	if err != nil {
		return nil, false
	}
	key := redis.GenerateCommentKey(vid)

	// 这里的user要实时查询最新的信息（因为关注、粉丝等信息是会动态更新的）
	user, ok := GetUserByID(comment.UserID, uid)
	// 这里获取用户失败，但是不影响新增评论的操作，所以只需要从redis中把数据删除即可
	if !ok {
		go func() {
			ctx, cancel := redis.WithTimeoutContextBySecond(2)
			defer cancel()
			redis.RdbComment.Del(ctx, key)
		}()
		return comment, true
	}
	// 如果user返回成功，将对应User的最新信息（关注、粉丝等）同步到评论中
	comment.User = *user

	// 创建一个用于接收error的通道
	errCh := make(chan error)

	// 并发更新redis中的对应视频vid的list
	go func() {
		ctx, cancel := redis.WithTimeoutContextBySecond(2)
		defer cancel()

		// 如果redis中已经有相关的key value，直接更新
		// 如果没有，则需要从数据库中加载
		if n, err := redis.RdbComment.Exists(ctx, key).Result(); n == 0 {
			if err != nil {
				errCh <- err // 发送错误信息到通道
				return
			}
			if _, err := updateRedisComments(ctx, key, vid, uid); err != nil {
				errCh <- err // 发送错误信息到通道
				return
			}
		}

		// 以下过程如果出现失败情况，都要把当前redis的对应视频评论数据清空，避免出现读取错误数据的情况
		commentJson, err := json.Marshal(comment)
		if err != nil {
			errCh <- err // 发送错误信息到通道
			return
		}

		if err := redis.RdbComment.HSet(ctx, key, comment.ID, commentJson).Err(); err != nil {
			errCh <- err // 发送错误信息到通道
			return
		}

		errCh <- nil // 发送nil表示没有错误
	}()

	// 如果更新redis失败，直接删除这个视频vid对应的hash结构
	go func() {
		ctx, cancel := redis.WithTimeoutContextBySecond(2)
		defer cancel()

		err = <-errCh
		if err != nil {
			redis.RdbComment.Del(ctx, key)
		}
	}()

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
	key := redis.GenerateCommentKey(vid)

	// 并发更新redis中的对应视频vid的list
	go func() {
		ctx, cancel := redis.WithTimeoutContextBySecond(2)
		defer cancel()

		// 如果redis中已经有相关的key value，直接更新
		// 如果没有，则需要从数据库中加载
		if n, err := redis.RdbComment.Exists(ctx, key).Result(); n == 0 {
			if err != nil {
				return
			}
			if _, err := updateRedisComments(ctx, key, vid, uid); err != nil {
				return
			}
		}

		// 如果要删除的数据同步到redis中执行删除失败，直接把当前视频的缓存清空，避免读取错误数据
		if err := redis.RdbComment.HDel(ctx, key, strconv.Itoa(int(commentID))).Err(); err != nil {
			redis.RdbComment.Del(ctx, key)
			return
		}
	}()

	// 异步删除评论，如果放到MQ失败则返回false，并回退redis的更改
	if rabbitmq.RMQDelComment.Publish(rabbitmq.GenerateDelCommentMQParam(commentID)) != nil {
		go func() {
			ctx, cancel := redis.WithTimeoutContextBySecond(2)
			defer cancel()

			redis.RdbComment.Del(ctx, key)
		}()
		return false
	}

	return true
}
