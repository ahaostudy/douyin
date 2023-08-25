package redis

import (
	"fmt"
	"main/config"
)

// GenerateMessageKey 通过 user_a_id 与 user_b_id 生成两个 key ，保证小 id 在前，大 id 在后
func GenerateMessageKey(aid, bid uint) (string, string) {
	var minID, maxID uint
	if aid < bid {
		minID, maxID = aid, bid
	} else {
		minID, maxID = bid, aid
	}
	minKey := fmt.Sprintf("%s:%d", config.RedisKeyOfMessage, minID)
	maxKey := fmt.Sprintf("%s:%d", config.RedisKeyOfMessage, maxID)
	return minKey, maxKey
}

// GenerateLikeKey 通过 video_id 生成 key
func GenerateLikeKey(id uint) string {
	return fmt.Sprintf("%s:%d", config.RedisKeyOfLike, id)
}
