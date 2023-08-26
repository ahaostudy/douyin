package service

import (
	"main/middleware/rabbitmq"
	"main/middleware/redis"
)

// FavoriteAction 赞业务
// uid, vid 分别对应用户ID和视频ID
// t 表示操作类型
func FavoriteAction(uid, vid uint, t int) bool {
	switch t {
	case 1:
		return like(uid, vid)
	case 2:
		return unLike(uid, vid)
	default:
		return false
	}
}

// 点赞业务
func like(uid, vid uint) bool {
	key := redis.GenerateLikeKey(vid)
	ctx, cancel := redis.WithTimeoutContextBySecond(3)
	defer cancel()

	// 1. 加载用户的喜欢列表和视频的点赞列表到redis
	if LoadLikeList(ctx, vid) != nil {
		return false
	}

	// 2. 更新redis数据，进行点赞操作
	//    如果已经有点赞记录或更新失败则返回false
	if n, err := redis.RdbLike.SAdd(ctx, key, uid).Result(); n == 0 || err != nil {
		return false
	}

	// 3. 写入redis成功，使用MQ异步更新MySQL
	if err := rabbitmq.RMQLike.Publish(rabbitmq.GenerateLikeMQParam(uid, vid)); err != nil {
		// 如果消息发送到MQ失败，则撤销redis的数据
		redis.RdbLike.SRem(ctx, key, uid)
		return false
	}

	// 4. 维护用户信息中的喜欢数和获赞数
	go func() {
		// 维护喜欢数
		// 判断redis是否有该用户信息缓存，无缓存时不需要更新
		if !ExistsUserInfo(ctx, uid) {
			return
		}
		key := redis.GenerateUserKey(uid)
		// 喜欢数+1，如果更新失败则删除key，保证数据一致
		if err := redis.RdbUser.HIncrBy(ctx, key, "favorite_count", 1).Err(); err != nil {
			redis.RdbUser.Del(ctx, key)
			return
		}

		// 维护获赞数
		authorID, err := GetAuthorID(ctx, vid)
		if err != nil || !ExistsUserInfo(ctx, authorID) {
			return
		}
		key = redis.GenerateUserKey(authorID)
		if err := redis.RdbUser.HIncrBy(ctx, key, "total_favorited", 1).Err(); err != nil {
			redis.RdbUser.Del(ctx, key)
			return
		}
	}()

	return true
}

// 取消点赞
// 取消点赞逻辑与点赞逻辑类似，不重复写注释
func unLike(uid, vid uint) bool {
	key := redis.GenerateLikeKey(vid)
	ctx, cancel := redis.WithTimeoutContextBySecond(3)
	defer cancel()

	if LoadLikeList(ctx, vid) != nil {
		return false
	}

	if err := redis.RdbLike.SRem(ctx, key, uid).Err(); err != nil {
		return false
	}

	if err := rabbitmq.RMQUnLike.Publish(rabbitmq.GenerateUnLikeMQParam(uid, vid)); err != nil {
		redis.RdbLike.SAdd(ctx, key, uid)
		return false
	}

	go func() {
		key := redis.GenerateUserKey(uid)
		if !ExistsUserInfo(ctx, uid) {
			return
		}
		if err := redis.RdbUser.HIncrBy(ctx, key, "favorite_count", -1).Err(); err != nil {
			redis.RdbUser.Del(ctx, key)
			return
		}

		authorID, err := GetAuthorID(ctx, vid)
		if err != nil || !ExistsUserInfo(ctx, authorID) {
			return
		}
		key = redis.GenerateUserKey(authorID)
		if err := redis.RdbUser.HIncrBy(ctx, key, "total_favorited", -1).Err(); err != nil {
			redis.RdbUser.Del(ctx, key)
			return
		}
	}()

	return true
}
