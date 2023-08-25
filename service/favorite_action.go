package service

import (
	"main/config"
	"main/dao"
	"main/middleware/rabbitmq"
	"main/middleware/redis"
	"main/model"
	"sync"
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
// TODO: 使用MQ异步更新MySQL，怎么保证Redis和MySQL数据一致
func like(uid, vid uint) bool {
	key := redis.GenerateLikeKey(vid)
	ctx, cancel := redis.WithTimeoutContextBySecond(3)
	defer cancel()

	// 1. 判断redis中是否有视频的点赞用户列表
	//    不存在时先从MySQL读取列表到redis
	if n, err := redis.RdbLike.Exists(ctx, key).Result(); n == 0 || err != nil {
		// 若读取失败，将该key从redis中删除，防止读取到一半，有数据但不全的情况
		if loadLikeListByVideo(vid) != nil {
			redis.RdbLike.Del(ctx, key)
			return false
		}
	}

	// 2. 更新redis数据，进行点赞操作
	//    如果已经有点赞记录或更新失败则返回false
	//    更新成功后刷新key过期时间，保证热点key不过期
	if n, err := redis.RdbLike.SAdd(ctx, key, uid).Result(); n == 0 || err != nil {
		return false
	}
	redis.RdbLike.Expire(ctx, key, config.RedisKeyTTL)

	// 3. 写入redis成功，使用MQ异步更新MySQL
	if err := rabbitmq.RMQLikeAdd.Publish(rabbitmq.GenerateLikeAddMQParam(uid, vid)); err != nil {
		// 如果消息发送到MQ失败，则撤销redis的数据
		redis.RdbLike.SRem(ctx, key, uid)
		return false
	}

	return true
}

// 取消点赞
// 取消点赞逻辑与点赞逻辑类似，不重复写注释
func unLike(uid, vid uint) bool {
	key := redis.GenerateLikeKey(vid)
	ctx, cancel := redis.WithTimeoutContextBySecond(3)
	defer cancel()

	if n, err := redis.RdbLike.Exists(ctx, key).Result(); n == 0 || err != nil {
		if loadLikeListByVideo(vid) != nil {
			redis.RdbLike.Del(ctx, key)
			return false
		}
	}

	if n, err := redis.RdbLike.SRem(ctx, key, uid).Result(); n == 0 || err != nil {
		return false
	}
	redis.RdbLike.Expire(ctx, key, config.RedisKeyTTL)

	if err := rabbitmq.RMQLikeDel.Publish(rabbitmq.GenerateLikeDelMQParam(uid, vid)); err != nil {
		redis.RdbLike.SAdd(ctx, key, uid)
		return false
	}

	return true
}

// 读取视频点赞列表到Redis
func loadLikeListByVideo(vid uint) error {
	// 从MySQL中获取数据
	likeList, err := dao.GetLikeListByVideoID(vid)
	if err != nil {
		return err
	}
	key := redis.GenerateLikeKey(vid)

	ctx, cancel := redis.WithTimeoutContextBySecond(2)
	defer cancel()

	// 先写入一个空value，这样就算key没有数据，也会保留在redis中，防止缓存穿透
	redis.RdbLike.SAdd(ctx, key, config.RedisValueOfNULL)
	// 设置过期时间
	redis.RdbLike.Expire(ctx, key, config.RedisKeyTTL)

	// 遍历点赞列表并发添加到Redis中
	wg, err := sync.WaitGroup{}, nil
	for _, like := range likeList {
		wg.Add(1)
		go func(like model.Like) {
			defer wg.Done()
			e := redis.RdbLike.SAdd(ctx, key, like.UserID).Err()
			if e != nil {
				err = e
			}
		}(*like)
	}
	wg.Wait()

	// 只要有一次写入发生错误，都要返回错误
	return err
}
