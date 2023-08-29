package service

import (
	"fmt"
	"main/config"
	"main/dao"
	"main/middleware/redis"
	"main/model"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// GetMessageList 获取消息列表
func GetMessageList(userID, toUserID uint, preMsgTime time.Time) ([]*model.Message, bool) {
	var messageList []*model.Message
	minKey, maxKey := redis.GenerateMessageKey(userID, toUserID)

	ctx, cancel := redis.WithTimeoutContextBySecond(3)
	defer cancel()

	// 0. 处理前端bug
	//    如果管道不为空，将当前请求时间戳传入该管道
	//    不处理前端bug时将此处删掉即可
	if userLatestMessageTime[userID] != nil {
		userLatestMessageTime[userID] <- preMsgTime.Truncate(time.Second)
	}

	// 1. 判断redis中是否用两位用户的消息缓存
	//    规定：主key（redis中的key）存储较小的ID，副key（redis中的hash的key）存储较大的ID
	latestTimeStr, err := redis.RdbMessage.HGet(ctx, minKey, maxKey).Result()
	if len(latestTimeStr) == 0 || err != nil {
		// 不存在时从MySQL读取数据到Redis
		latestTime, err := dao.GetLatestMessageTime(userID, toUserID)
		if err == gorm.ErrRecordNotFound {
			latestTime = time.UnixMilli(0)
		} else if err != nil {
			return messageList, false
		}
		fmt.Println(latestTime)
		redis.RdbMessage.HSet(ctx, minKey, maxKey, latestTime.UnixMilli())
		latestTimeStr = strconv.FormatInt(latestTime.UnixMilli(), 10)
	}
	// 更新有效时间
	redis.RdbMessage.Expire(ctx, minKey, config.RedisKeyTTL)

	// 2. 解析时间戳
	latestTimeStamp, err := strconv.ParseInt(latestTimeStr, 10, 64)
	if err != nil {
		return messageList, false
	}
	latestTime := time.UnixMilli(latestTimeStamp)

	// 3. 比较两个时间
	//    如果请求的pre_msg_time与最新消息时间相同或在其之后，则跳过
	if !preMsgTime.Before(latestTime) {
		return messageList, true
	}

	// 4. pre_msg_time在最新消息时间之前，从数据库获取聊天记录返回
	messageList, err = dao.GetMessageList(userID, toUserID, preMsgTime)
	return messageList, err == nil
}
