package service

import (
	"main/middleware/rabbitmq"
	"main/middleware/redis"
)

// RelationAction 关注操作
// uid: 用户ID, tid: 对方ID, t: 操作类型
func RelationAction(uid, tid uint, t int) bool {
	switch t {
	case 1:
		return follow(uid, tid)
	case 2:
		return unFollow(uid, tid)
	default:
		return false
	}
}

// 关注
func follow(uid, tid uint) bool {
	ctx, cancel := redis.WithTimeoutContextBySecond(300)
	defer cancel()

	// 加载关注和粉丝列表到redis
	if LoadFollowList(ctx, uid) != nil || LoadFollowerList(ctx, tid) != nil {
		return false
	}

	// 更新redis
	// 如果已经关注或更新失败则返回false
	// 如果更新第二个set失败，则回滚第一个set的操作
	fk, fek := redis.GenerateFollowKey(uid), redis.GenerateFollowerKey(tid)
	if n, err := redis.RdbFollow.SAdd(ctx, fk, tid).Result(); n == 0 || err != nil {
		return false
	}
	if n, err := redis.RdbFollow.SAdd(ctx, fek, uid).Result(); n == 0 || err != nil {
		redis.RdbFollow.SRem(ctx, fk, tid)
		return false
	}

	// 异步更新数据库
	// TODO：本地测试先注释掉涉及MQ的部分
	if rabbitmq.RMQFollow.Publish(rabbitmq.GenerateFollowMQParam(uid, tid)) != nil {
		redis.RdbFollow.SRem(ctx, fk, tid)
		redis.RdbFollow.SRem(ctx, fek, uid)
		return false
	}

	// 维护用户信息
	go func() {
		ctx, cancel := redis.WithTimeoutContextBySecond(2)
		defer cancel()

		// 更新关注数，如果更新失败则删除key，保证数据一致
		if ExistsUserInfo(ctx, uid) {
			key := redis.GenerateUserKey(uid)
			if redis.RdbUser.HIncrBy(ctx, key, "follow_count", 1).Err() != nil {
				redis.RdbUser.Del(ctx, key)
				return
			}
		}

		// 更新粉丝数
		if ExistsUserInfo(ctx, tid) {
			key := redis.GenerateUserKey(tid)
			if redis.RdbUser.HIncrBy(ctx, key, "follower_count", 1).Err() != nil {
				redis.RdbUser.Del(ctx, key)
				return
			}
		}
	}()

	return true
}

// 取消关注
func unFollow(uid, tid uint) bool {
	ctx, cancel := redis.WithTimeoutContextBySecond(2)
	defer cancel()

	if LoadFollowList(ctx, uid) != nil || LoadFollowerList(ctx, tid) != nil {
		return false
	}

	fk, fek := redis.GenerateFollowKey(uid), redis.GenerateFollowerKey(tid)
	if n, err := redis.RdbFollow.SRem(ctx, fk, tid).Result(); n == 0 || err != nil {
		return false
	}
	if n, err := redis.RdbFollow.SRem(ctx, fek, uid).Result(); n == 0 || err != nil {
		redis.RdbFollow.SAdd(ctx, fk, tid)
		return false
	}

	if rabbitmq.RMQUnFollow.Publish(rabbitmq.GenerateUnFollowMQParam(uid, tid)) != nil {
		redis.RdbFollow.SAdd(ctx, fk, tid)
		redis.RdbFollow.SAdd(ctx, fek, uid)
		return false
	}

	go func() {
		ctx, cancel := redis.WithTimeoutContextBySecond(2)
		defer cancel()

		if ExistsUserInfo(ctx, uid) {
			key := redis.GenerateUserKey(uid)
			if redis.RdbUser.HIncrBy(ctx, key, "follow_count", -1).Err() != nil {
				redis.RdbUser.Del(ctx, key)
				return
			}
		}

		if ExistsUserInfo(ctx, tid) {
			key := redis.GenerateUserKey(tid)
			if redis.RdbUser.HIncrBy(ctx, key, "follower_count", -1).Err() != nil {
				redis.RdbUser.Del(ctx, key)
				return
			}
		}
	}()

	return true
}
