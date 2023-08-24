package service

import (
	"main/dao"
	"main/middleware/redis"
	"main/model"
	"time"
)

var userLatestMessageTime = make(map[uint]chan time.Time)

// 发送消息业务（不处理前端bug的版本）
func insertMessage(fromID, toID uint, content string) bool {
	// 插入一条消息到数据库
	message, err := dao.InsertMessage(&model.Message{
		FromUserID: fromID,
		ToUserID:   toID,
		Content:    content,
	})
	if err != nil {
		return false
	}

	minKey, maxKey := generateMessageKey(fromID, toID)

	ctx, cancel := redis.WithTimeoutContextBySecond(2)
	defer cancel()

	// 更新redis缓存的最新消息时间
	// 更新失败时删除该缓存，保证数据一致性
	// TODO: 需要加锁，防止新消息更新到redis的速度比旧消息快，造成旧消息的时间覆盖新消息
	if err := redis.RdbMessage.HSet(ctx, minKey, maxKey, message.CreatedAt.UnixMilli()).Err(); err != nil {
		redis.RdbMessage.HDel(ctx, minKey, maxKey)
	}

	return true
}

// InsertMessage 发送消息业务
func InsertMessage(fromID, toID uint, content string) bool {
	// 创建一个管道，提供给chat接口
	// chat接口会将请求时间戳传入，作为当前这条消息的发送时间
	if userLatestMessageTime[fromID] != nil {
		return false
	}
	userLatestMessageTime[fromID] = make(chan time.Time)

	// 新建一个协程来等待chat接口传递参数并执行业务
	go func() {
		// 阻塞等待chat接口传递时间戳
		t := <-userLatestMessageTime[fromID]
		// 接收完后删除时间戳
		userLatestMessageTime[fromID] = nil
		// 插入一条消息到数据库
		message, err := dao.InsertMessage(&model.Message{
			FromUserID: fromID,
			ToUserID:   toID,
			Content:    content,
			CreatedAt:  t,
		})
		if err != nil {
			return
		}

		minKey, maxKey := generateMessageKey(fromID, toID)

		ctx, cancel := redis.WithTimeoutContextBySecond(2)
		defer cancel()

		// 更新redis缓存的最新消息时间
		// 更新失败时删除该缓存，保证数据一致性
		// TODO: 需要加锁，防止新消息更新到redis的速度比旧消息快，造成旧消息的时间覆盖新消息
		if err := redis.RdbMessage.HSet(ctx, minKey, maxKey, message.CreatedAt.UnixMilli()).Err(); err != nil {
			redis.RdbMessage.HDel(ctx, minKey, maxKey)
		}
	}()

	return true
}
