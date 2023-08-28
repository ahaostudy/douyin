package redis

import (
	"fmt"
	"main/config"
)

func GenerateLikeKey(id uint) string {
	return fmt.Sprintf("%s:%d", config.RedisKeyOfLike, id)
}

func GenerateFollowKey(id uint) string {
	return fmt.Sprintf("%s:%d", config.RedisKeyOfFollow, id)
}

func GenerateFollowerKey(id uint) string {
	return fmt.Sprintf("%s:%d", config.RedisKeyOfFollower, id)
}

func GenerateUserKey(id uint) string {
	return fmt.Sprintf("%s:%d", config.RedisKeyOfUser, id)
}

func GenerateAuthorKey(id uint) string {
	return fmt.Sprintf("%s:%d", config.RedisKeyOfAuthor, id)
}

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

// 评论是依赖于视频的，所以这里redis 主key是视频vid，副key是视频的各个评论的
func GenerateCommentKey(vid uint) string {
	return fmt.Sprintf("%s:%d", config.RedisKeyOfComment, vid)
}
