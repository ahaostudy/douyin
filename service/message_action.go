package service

import (
	"fmt"
	"main/dao"
	"main/middleware/redis"
	"main/model"
	"strconv"
	"time"
)

var userLatestMessageTime = make(map[uint]chan time.Time)

// 发送消息业务（不处理前端bug的版本）
func sentMessage(fromID, toID uint, content string) bool {
	// 插入一条消息到数据库
	message, err := dao.InsertMessage(&model.Message{
		FromUserID: fromID,
		ToUserID:   toID,
		Content:    content,
	})
	if err != nil {
		return false
	}

	// 并发更新redis缓存的最新消息时间
	// 更新失败时删除该缓存，保证数据一致性
	// 使用乐观锁，防止新消息更新到redis的速度比旧消息快，造成旧消息的时间覆盖新消息
	go func() {
		minKey, maxKey := redis.GenerateMessageKey(fromID, toID)

		ctx, cancel := redis.WithTimeoutContextBySecond(2)
		defer cancel()

		// 获取乐观锁，保证并发安全
		lockKey := generateLockKey(minKey, maxKey)
		lockID, err := redis.Lock(redis.RdbMessage, lockKey)
		defer redis.Unlock(redis.RdbMessage, lockKey, lockID)
		if err != nil {
			redis.RdbMessage.HDel(ctx, minKey, maxKey)
			return
		}

		// 获取缓存中的时间戳，判断该时间戳与将要写入的时间戳哪个更新
		cacheStampStr, err := redis.RdbMessage.HGet(ctx, minKey, maxKey).Result()
		if err != nil || len(cacheStampStr) == 0 {
			return
		}
		cacheStamp, err := strconv.ParseInt(cacheStampStr, 10, 64)
		if err != nil {
			redis.RdbMessage.HDel(ctx, minKey, maxKey)
			return
		}

		// 如果缓存的时间戳在将要写入的时间戳之前，则更新时间戳
		if time.UnixMilli(cacheStamp).Before(message.CreatedAt) {
			if err := redis.RdbMessage.HSet(ctx, minKey, maxKey, message.CreatedAt.UnixMilli()).Err(); err != nil {
				redis.RdbMessage.HDel(ctx, minKey, maxKey)
			}
		}
	}()

	return true
}

// SendMessage 发送消息业务
func SendMessage(fromID, toID uint, content string) bool {
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
		// 接收完后删除管道
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

		minKey, maxKey := redis.GenerateMessageKey(fromID, toID)

		ctx, cancel := redis.WithTimeoutContextBySecond(2)
		defer cancel()

		// 更新redis缓存的最新消息时间
		// 更新失败时删除该缓存，保证数据一致性
		// 使用乐观锁，防止新消息更新到redis的速度比旧消息快，造成旧消息的时间覆盖新消息
		lockKey := generateLockKey(minKey, maxKey)
		lockID, err := redis.Lock(redis.RdbMessage, lockKey)
		defer redis.Unlock(redis.RdbMessage, lockKey, lockID)
		if err != nil {
			redis.RdbMessage.HDel(ctx, minKey, maxKey)
			return
		}

		// 获取缓存中的时间戳，判断该时间戳与将要写入的时间戳哪个更新
		cacheStampStr, err := redis.RdbMessage.HGet(ctx, minKey, maxKey).Result()
		if err != nil || len(cacheStampStr) == 0 {
			return
		}
		cacheStamp, err := strconv.ParseInt(cacheStampStr, 10, 64)
		if err != nil {
			redis.RdbMessage.HDel(ctx, minKey, maxKey)
			return
		}

		// 如果缓存的时间戳在将要写入的时间戳之前，则更新时间戳
		if time.UnixMilli(cacheStamp).Before(message.CreatedAt) {
			if err := redis.RdbMessage.HSet(ctx, minKey, maxKey, message.CreatedAt.UnixMilli()).Err(); err != nil {
				redis.RdbMessage.HDel(ctx, minKey, maxKey)
			}
		}
	}()

	return true
}

// 生成lock的key
func generateLockKey(minKey, maxKey string) string {
	return fmt.Sprintf("%s-%s", minKey, maxKey)
}
