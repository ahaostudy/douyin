package service

import (
	"main/config"
	"main/middleware/redis"
	"main/model"
	"strconv"
	"sync"
)

type FriendPair struct {
	Idx    int
	Friend *model.User
}

func GetFriendList(id, curID uint) ([]*model.User, bool) {
	ctx, cancel := redis.WithTimeoutContextBySecond(2)
	defer cancel()
	followKey, followerKey := redis.GenerateFollowKey(id), redis.GenerateFollowerKey(id)

	// 1. 加载关注列表和粉丝列表到redis
	wg, status := sync.WaitGroup{}, true
	wg.Add(2)
	go func() {
		defer wg.Done()
		if LoadFollowList(ctx, id) != nil {
			status = false
		}
	}()
	go func() {
		defer wg.Done()
		if LoadFollowerList(ctx, id) != nil {
			status = false
		}
	}()
	wg.Wait()
	if !status {
		return nil, false
	}

	// 2. 从redis中获取好友列表
	friends, err := redis.RdbFollow.SInter(ctx, followKey, followerKey).Result()
	if err != nil {
		return nil, false
	}

	// 3. 根据朋友列表的用户ID获取用户信息
	friendChan := make(chan FriendPair)
	status, nullIdx := true, len(friends)
	for i, friend := range friends {
		go func(i int, friend string) {
			if friend == config.RedisValueOfNULL {
				nullIdx = i
				return
			}
			_fid, _ := strconv.ParseUint(friend, 10, 64)
			fid := uint(_fid)
			user, ok := GetUserByID(fid, curID)
			if !ok {
				status = false
			}
			friendChan <- FriendPair{Idx: i, Friend: user}
		}(i, friend)
	}

	// 4. 通过管道获取结果，添加到结果切片中
	if len(friends) <= 1 {
		return nil, len(friends) == 1
	}
	friendList := make([]*model.User, len(friends)-1)
	for i := 1; i < len(friends); i++ {
		friend := <-friendChan
		idx := friend.Idx
		if idx >= nullIdx {
			idx--
		}
		friendList[idx] = friend.Friend
	}

	return friendList, status
}
