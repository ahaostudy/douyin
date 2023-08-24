package service

import (
	"fmt"
	"gorm.io/gorm"
	"main/config"
	"main/dao"
	"main/middleware/redis"
	"main/model"
	"strconv"
	"time"
)

// GetMessageList 获取消息列表
func GetMessageList(userID, toUserID uint, preMsgTime time.Time) ([]*model.Message, bool) {
	var messageList []*model.Message
	minKey, maxKey := generateMessageKey(userID, toUserID)

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
		if err != nil && err != gorm.ErrRecordNotFound {
			return messageList, false
		}
		fmt.Println(latestTime.UnixMilli())
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

// TODO: 这些 generate key 的函数感觉统一起来管理更好
// 通过 user_a_id 与 user_b_id 生成两个 key ，保证小 id 在前，大 id 在后
func generateMessageKey(a, b uint) (string, string) {
	var minID, maxID uint
	if a < b {
		minID, maxID = a, b
	} else {
		minID, maxID = b, a
	}
	minKey := fmt.Sprintf("%s:%d", config.RedisKeyOfMessage, minID)
	maxKey := fmt.Sprintf("%s:%d", config.RedisKeyOfMessage, maxID)
	return minKey, maxKey
}
